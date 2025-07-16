package handlers

import (
	"Auth-service/internal/service"
	proto "Auth-service/pkg/proto/auth"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	service service.AuthServiceServer
}

func AuthHandlerConstructor(
	service service.AuthServiceServer,
) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Registration(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
	user, err := h.service.Registrations(ctx, req.FullName, req.Email, req.Password, req.CurrencyCode)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.RegistrationResponse{
		UserID:  user.ID,
		Message: "Created user successfully!",
	}, nil
}

func (h *AuthHandler) Authorization(ctx context.Context, req *proto.AuthorizationRequest) (*proto.AuthorizationResponse, error) {
	user, token, err := h.service.Authorizations(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.AuthorizationResponse{
		UserID:  user.ID,
		Token:   token,
		Message: "Authorization successfully!",
	}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *proto.ChangePasswordRequest) (*proto.ChangePasswordResponse, error) {
	err := h.service.ChangePasswords(ctx, req.Email, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.ChangePasswordResponse{
		Message: "Change password successfully!",
	}, nil
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
