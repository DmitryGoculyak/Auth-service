package service

import (
	cfg "Auth-service/internal/config"
	"Auth-service/internal/dto"
	"Auth-service/pkg/client/billing"
	"Auth-service/pkg/client/currency"
	"Auth-service/pkg/jwt"
	proto "Auth-service/pkg/proto/auth"
	"Auth-service/pkg/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"
	"log"
	"net"
	"time"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	billing   *billing.BillingClient
	currency  *currency.CurrencyClient
	validator *validator.Validate
	db        *sqlx.DB
}

func AuthServerConstructor(
	db *sqlx.DB,
	billingClient *billing.BillingClient,
	currencyClient *currency.CurrencyClient,
	validator *validator.Validate,
) *AuthServer {
	return &AuthServer{
		db:        db,
		billing:   billingClient,
		currency:  currencyClient,
		validator: validator,
	}
}

func (s *AuthServer) Registration(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {

	input := dto.RegistrationInput{
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
		Password:     req.GetPassword(),
		CurrencyCode: req.GetCurrencyCode(),
	}

	if err := s.validator.Struct(input); err != nil {
		return nil, status.Error(codes.Internal, "Validation error: "+err.Error())
	}

	tx := s.db.MustBegin()
	var err error

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			log.Printf("Panic: %v", r)
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var userID string
	err = tx.Get(&userID, "INSERT INTO users(full_name) VALUES($1) RETURNING id;", req.FullName)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO emails (user_id,email) VALUES($1, $2)", userID, req.Email)
	if err != nil {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO passwords(user_id, hash) VALUES($1, $2)", userID, hash)
	if err != nil {
		return nil, err
	}

	if err = s.currency.CheckCurrencyExists(ctx, req.CurrencyCode); err != nil {
		return nil, err
	}

	if err = s.billing.CreateWallet(ctx, userID, req.CurrencyCode); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &proto.RegistrationResponse{
		UserID:  userID,
		Message: "Registered successfully!",
	}, nil
}

func (s *AuthServer) Authorization(ctx context.Context, req *proto.AuthorizationRequest) (*proto.AuthorizationResponse, error) {
	var userID, hash string

	err := s.db.QueryRowContext(ctx, `
		SELECT u.id, p.hash FROM users u
		JOIN emails e ON u.id = e.user_id
		JOIN passwords p ON u.id = p.user_id
		WHERE e.email = $1
	`, req.Email).Scan(&userID, &hash)

	if err != nil {
		log.Printf("[AUTH] Error fetching user and hash for email %s: %v", req.Email, err)
		return nil, status.Error(codes.PermissionDenied, "invalid email or user not found")
	}

	if !utils.CheckPassword(hash, req.Password) {
		log.Printf("[AUTH] Password mismatch error for userID %s: %v", userID, err)
		return nil, status.Error(codes.PermissionDenied, "invalid password")
	}

	token, err := jwt.GenerateJWTToken(userID)
	if err != nil {
		log.Printf("[AUTH] Failed to generate JWT token for userID %s: %v", userID, err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &proto.AuthorizationResponse{
		UserID:  userID,
		Token:   token,
		Message: "Authorized successfully!",
	}, nil
}

func (s *AuthServer) DeleteAllUsers(ctx context.Context, _ *proto.Empty) (*proto.DeleteResponse, error) {
	_, err := s.db.Exec("DELETE FROM users;")
	if err != nil {
		return nil, err
	}
	return &proto.DeleteResponse{Message: "All users deleted"}, nil
}

func RunServer(cfg *cfg.GrpcServiceConfig, server *AuthServer) {

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterAuthServiceServer(s, server)

	log.Printf("[gRPC] Server started at time %v on address %v",
		time.Now().Format("[2006-01-02] [15:04]"), address)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
