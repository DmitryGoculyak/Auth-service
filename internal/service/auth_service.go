package service

import (
	"Auth-service/internal/dto"
	"Auth-service/internal/repository"
	"Auth-service/pkg/client/billing"
	"Auth-service/pkg/client/currency"
	"Auth-service/pkg/jwt"
	proto "Auth-service/pkg/proto/auth"
	"Auth-service/pkg/utils"

	"context"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer interface {
	RegistrationUser(ctx context.Context, input *dto.RegistrationInput) (*proto.RegistrationResponse, error)
	AuthorizationUser(ctx context.Context, input *dto.AuthorizationInput) (*proto.AuthorizationResponse, error)
	DeleteUsers(ctx context.Context) error
}

type AuthService struct {
	db       *sqlx.DB
	repo     repository.AuthRepository
	billing  *billing.BillingClient
	currency *currency.CurrencyClient
}

func AuthServiceConstructor(
	db *sqlx.DB,
	repo repository.AuthRepository,
	billingClient *billing.BillingClient,
	currencyClient *currency.CurrencyClient,
) *AuthService {
	return &AuthService{
		db:       db,
		billing:  billingClient,
		currency: currencyClient,
		repo:     repo,
	}
}

func (s *AuthService) RegistrationUser(ctx context.Context, input *dto.RegistrationInput) (*proto.RegistrationResponse, error) {
	tx := s.db.MustBegin()
	var err error

	defer func() {
		if rec := recover(); rec != nil {
			_ = tx.Rollback()
			log.Printf("Panic: %v", rec)
			panic(rec)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	userID, err := s.repo.CreateUser(tx, input.FullName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = s.repo.SaveEmail(tx, userID, input.Email); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	if err = s.repo.SavePassword(tx, userID, hash); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = s.billing.CreateWallet(ctx, userID, input.CurrencyCode); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = s.currency.CheckCurrencyExists(ctx, input.CurrencyCode); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.RegistrationResponse{
		UserID:  userID,
		Message: "User registered",
	}, nil
}

func (s *AuthService) AuthorizationUser(ctx context.Context, input *dto.AuthorizationInput) (*proto.AuthorizationResponse, error) {
	userID, hash, err := s.repo.FindUserEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(hash, input.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateJWTToken(userID)
	if err != nil {
		return nil, err
	}

	return &proto.AuthorizationResponse{
		UserID:  userID,
		Token:   token,
		Message: "User authorized",
	}, nil
}

func (s *AuthService) DeleteUsers(ctx context.Context) error {
	return s.repo.DeleteAllUsers(ctx)
}
