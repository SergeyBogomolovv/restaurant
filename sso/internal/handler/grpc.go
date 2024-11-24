package handler

import (
	"context"
	"time"

	pb "github.com/SergeyBogomolovv/restaurant/common/api/gen/sso"
	er "github.com/SergeyBogomolovv/restaurant/common/errors"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUsecase interface{}

type RegisterUsecase interface {
	RegisterCustomer(ctx context.Context, dto *domain.RegisterCustomerDTO) (string, error)
	RegisterWaiter(ctx context.Context, dto *domain.RegisterWaiterDTO, token string) (string, error)
	RegisterAdmin(ctx context.Context, dto *domain.RegisterAdminDTO, token string) (string, error)
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
	dto := &domain.RegisterCustomerDTO{
		Email:     req.Email,
		Name:      req.Name,
		Birthdate: time.Unix(req.Birthdate, 0),
		Password:  []byte(req.Password),
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	entityID, err := h.register.RegisterCustomer(ctx, dto)
	if err != nil {
		if err == er.ErrCustomerAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, "Customer with this email already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to register customer, error: %v", err)
	}
	return &pb.RegisterResponse{EntityId: entityID}, nil
}

func (h *ssoHandler) RegisterWaiter(ctx context.Context, req *pb.RegisterWaiterRequest) (*pb.RegisterResponse, error) {
	dto := &domain.RegisterWaiterDTO{
		Login:     req.Login,
		Password:  []byte(req.Password),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Token:     req.SecretToken,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}

	entityID, err := h.register.RegisterWaiter(ctx, dto, req.SecretToken)
	if err != nil {
		//TODO: handle errors
		return nil, status.Errorf(codes.Internal, "failed to register waiter, error: %v", err)
	}
	return &pb.RegisterResponse{EntityId: entityID}, nil
}

func (h *ssoHandler) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest) (*pb.RegisterResponse, error) {
	dto := &domain.RegisterAdminDTO{
		Note:     req.Note,
		Login:    req.Login,
		Password: []byte(req.Password),
		Token:    req.SecretToken,
	}
	if err := h.validate.Struct(dto); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload, error: %v", err)
	}
	entityID, err := h.register.RegisterAdmin(ctx, dto, req.SecretToken)
	if err != nil {
		//TODO: handle errors
		return nil, status.Errorf(codes.Internal, "failed to register admin, error: %v", err)
	}
	return &pb.RegisterResponse{EntityId: entityID}, nil
}

func (h *ssoHandler) LoginCustomer(ctx context.Context, req *pb.LoginCustomerRequest) (*pb.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginCustomer not implemented")
}
func (h *ssoHandler) LoginWaiter(ctx context.Context, req *pb.LoginEmployeeRequest) (*pb.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginWaiter not implemented")
}
func (h *ssoHandler) LoginAdmin(ctx context.Context, req *pb.LoginEmployeeRequest) (*pb.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginAdmin not implemented")
}
func (h *ssoHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (h *ssoHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
