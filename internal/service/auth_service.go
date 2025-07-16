package service

import (
	"Auth-service/internal/dto"
	"Auth-service/internal/entity"
	"Auth-service/internal/repository"
	"Auth-service/pkg/client/billing"
	"Auth-service/pkg/client/currency"
	"Auth-service/pkg/jwt"
	"Auth-service/pkg/utils"
	"context"
	"github.com/go-playground/validator/v10"
	"log"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer interface {
	Registrations(ctx context.Context, fullName, email, password, currencyCode string) (*entity.User, error)
	Authorizations(ctx context.Context, email, password string) (*entity.User, string, error)
	DeleteUsers(ctx context.Context) error
}

type AuthService struct {
	db       *sqlx.DB
	repo     repository.AuthRepository
	validate *validator.Validate
	jwtToken *jwt.JWTConfig
	billing  *billing.BillingClient
	currency *currency.CurrencyClient
}

func AuthServiceConstructor(
	db *sqlx.DB,
	repo repository.AuthRepository,
	validate *validator.Validate,
	jwtToken *jwt.JWTConfig,
	billingClient *billing.BillingClient,
	currencyClient *currency.CurrencyClient,
) *AuthService {
	return &AuthService{
		db:       db,
		repo:     repo,
		validate: validate,
		jwtToken: jwtToken,
		billing:  billingClient,
		currency: currencyClient,
	}
}

func (s *AuthService) Registrations(ctx context.Context, fullName, email, password, currencyCode string) (u *entity.User, err error) {
	input := dto.RegistrationInput{
		FullName:     fullName,
		Email:        email,
		Password:     password,
		CurrencyCode: currencyCode,
	}

	if err = s.validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hash, err := utils.CreateHash(input.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	defer func() {
		if rec := recover(); rec != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[ERROR] Rollback error: %v", rbErr)
			}
			panic(rec)
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[ERROR] Rollback error: %v", rbErr)
			}
			return
		}
		if cmErr := tx.Commit(); cmErr != nil {
			log.Printf("[ERROR] Commit error: %v", cmErr)
			err = status.Error(codes.Internal, "commit failed: "+cmErr.Error())
		}
	}()

	user := &entity.User{
		FullName: input.FullName,
	}

	createUser, err := s.repo.CreateUser(tx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userEmail := &entity.UserEmail{
		UserID: createUser.ID,
		Email:  input.Email,
	}

	_, err = s.repo.SaveEmail(tx, userEmail)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userPassword := &entity.UserPassword{
		UserID:   createUser.ID,
		Password: hash,
	}
	_, err = s.repo.SavePassword(tx, userPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = s.currency.CheckCurrencyExists(ctx, input.CurrencyCode); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = s.billing.CreateWallet(ctx, user.ID, input.CurrencyCode); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return createUser, nil
}

func (s *AuthService) Authorizations(ctx context.Context, email, password string) (*entity.User, string, error) {
	input := dto.AuthorizationInput{
		Email:    email,
		Password: password,
	}

	if err := s.validate.Struct(input); err != nil {
		return nil, "", status.Error(codes.InvalidArgument, err.Error())
	}

	userId, hash, err := s.repo.FindUserEmail(ctx, input.Email)
	if err != nil {
		return nil, "", status.Error(codes.InvalidArgument, err.Error())
	}

	if !utils.CheckPassword(hash, input.Password) {
		return nil, "", status.Error(codes.FailedPrecondition, "password is incorrect")
	}

	token, err := s.jwtToken.GenerateToken(userId)
	if err != nil {
		return nil, "", status.Error(codes.InvalidArgument, err.Error())
	}

	return &entity.User{
		ID: userId,
	}, token, nil
}

func (s *AuthService) DeleteUsers(ctx context.Context) error {
	return s.repo.DeleteAllUsers(ctx)
}
