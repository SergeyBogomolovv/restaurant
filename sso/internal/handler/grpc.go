package handler

import (
	"context"
	"errors"
	"time"

	pb "github.com/SergeyBogomolovv/restaurant/common/api/gen/sso"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUsecase interface {
	LoginCustomer(ctx context.Context, dto *dto.LoginCustomerDTO) (*dto.TokensDTO, error)
	LoginWaiter(ctx context.Context, dto *dto.LoginEmployeeDTO) (*dto.TokensDTO, error)
	LoginAdmin(ctx context.Context, dto *dto.LoginEmployeeDTO) (*dto.TokensDTO, error)
	Refresh(ctx context.Context, token string) (string, error)
	Logout(ctx context.Context, token string) error
}

type RegisterUsecase interface {
	RegisterCustomer(ctx context.Context, dto *dto.RegisterCustomerDTO) (uuid.UUID, error)
	RegisterWaiter(ctx context.Context, dto *dto.RegisterWaiterDTO) (uuid.UUID, error)
	RegisterAdmin(ctx context.Context, dto *dto.RegisterAdminDTO) (uuid.UUID, error)
}

type ssoHandler struct {
	validate *validator.Validate
	auth     AuthUsecase
	register RegisterUsecase
	pb.UnimplementedSSOServer
}

func RegisterGRPCHandler(server *grpc.Server, auth AuthUsecase, register RegisterUsecase) {
	handler := &ssoHandler{
		validate: validator.New(validator.WithRequiredStructEnabled()),
		auth:     auth,
		register: register,
	}
	pb.RegisterSSOServer(server, handler)
}

func (h *ssoHandler) RegisterCustomer(ctx context.Context, req *pb.RegisterCustomerRequest) (*pb.RegisterResponse, error) {
	dto := &dto.RegisterCustomerDTO{
		Email:     req.Email,
		Name:      req.Name,
		Birthdate: time.Unix(req.Birthdate, 0),
		Password:  req.Password,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	entityID, err := h.register.RegisterCustomer(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrCustomerAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "Customer with this email already exists")
		default:
			return nil, status.Error(codes.Internal, "failed to register customer")
		}
	}
	return &pb.RegisterResponse{EntityId: entityID.String()}, nil
}

func (h *ssoHandler) RegisterWaiter(ctx context.Context, req *pb.RegisterWaiterRequest) (*pb.RegisterResponse, error) {
	dto := &dto.RegisterWaiterDTO{
		Login:     req.Login,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Token:     req.SecretToken,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}

	entityID, err := h.register.RegisterWaiter(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidSecretToken):
			return nil, status.Error(codes.Unauthenticated, "invalid secret token")
		case errors.Is(err, errs.ErrWaiterAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "Waiter with this login already exists")
		default:
			return nil, status.Error(codes.Internal, "failed to register waiter")
		}
	}
	return &pb.RegisterResponse{EntityId: entityID.String()}, nil
}

func (h *ssoHandler) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest) (*pb.RegisterResponse, error) {
	dto := &dto.RegisterAdminDTO{
		Note:     req.Note,
		Login:    req.Login,
		Password: req.Password,
		Token:    req.SecretToken,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	entityID, err := h.register.RegisterAdmin(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidSecretToken):
			return nil, status.Error(codes.Unauthenticated, "invalid secret token")
		case errors.Is(err, errs.ErrAdminAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "Admin with this login already exists")
		default:
			return nil, status.Error(codes.Internal, "failed to register admin")
		}
	}
	return &pb.RegisterResponse{EntityId: entityID.String()}, nil
}

func (h *ssoHandler) LoginCustomer(ctx context.Context, req *pb.LoginCustomerRequest) (*pb.LoginResponse, error) {
	dto := &dto.LoginCustomerDTO{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	tokens, err := h.auth.LoginCustomer(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Error(codes.Internal, "failed to login customer")
		}
	}
	return &pb.LoginResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (h *ssoHandler) LoginWaiter(ctx context.Context, req *pb.LoginEmployeeRequest) (*pb.LoginResponse, error) {
	dto := &dto.LoginEmployeeDTO{
		Login:    req.Login,
		Password: req.Password,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	tokens, err := h.auth.LoginWaiter(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Error(codes.Internal, "failed to login waiter")
		}
	}
	return &pb.LoginResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (h *ssoHandler) LoginAdmin(ctx context.Context, req *pb.LoginEmployeeRequest) (*pb.LoginResponse, error) {
	dto := &dto.LoginEmployeeDTO{
		Login:    req.Login,
		Password: req.Password,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	tokens, err := h.auth.LoginAdmin(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Error(codes.Internal, "failed to login admin")
		}
	}
	return &pb.LoginResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (h *ssoHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	token, err := h.auth.Refresh(ctx, req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidJwtToken):
			return nil, status.Error(codes.Unauthenticated, "invalid refreshToken")
		default:
			return nil, status.Error(codes.Internal, "failed to refresh token")
		}

	}
	return &pb.RefreshResponse{AccessToken: token}, nil
}

func (h *ssoHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if err := h.auth.Logout(ctx, req.RefreshToken); err != nil {
		return &pb.LogoutResponse{Status: "error"}, nil
	}
	return &pb.LogoutResponse{Status: "OK"}, nil
}
