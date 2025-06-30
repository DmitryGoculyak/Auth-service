package handlers

import (
	"Auth-service/internal/dto"
	"Auth-service/internal/service"
	proto "Auth-service/pkg/proto/auth"
	"context"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	service   service.AuthServiceServer
	validator *validator.Validate
}

func AuthHandlerConstructor(
	service service.AuthServiceServer,
	validator *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		service:   service,
		validator: validator,
	}
}

func (h *AuthHandler) Registration(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	input := dto.RegistrationInput{
		FullName:     req.FullName,
		Email:        req.Email,
		Password:     req.Password,
		CurrencyCode: req.CurrencyCode,
	}
	if err := h.validator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Registration validation failure")
	}

	return h.service.RegistrationUser(ctx, &input)
}

func (h *AuthHandler) Authorization(ctx context.Context, req *proto.AuthorizationRequest) (*proto.AuthorizationResponse, error) {
	input := dto.AuthorizationInput{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.validator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Authorization validation failure")
	}

	return h.service.AuthorizationUser(ctx, &input)
}

func (h *AuthHandler) DeleteAllUsers(ctx context.Context, _ *proto.Empty) (*proto.DeleteResponse, error) {
	err := h.service.DeleteUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Delete all users failure")
	}

	return &proto.DeleteResponse{
		Message: "All users deleted",
	}, nil
}
