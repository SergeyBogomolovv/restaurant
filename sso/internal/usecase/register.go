package usecase

import (
	"context"
	"log/slog"

	er "github.com/SergeyBogomolovv/restaurant/common/errors"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registerUsecase struct {
	customers CustomerRepo
	admins    AdminRepo
	waiters   WaiterRepo
	log       *slog.Logger
	secretKey string
}

func NewRegisterUsecase(log *slog.Logger, customers CustomerRepo, waiters WaiterRepo, admins AdminRepo, secretKey string) *registerUsecase {
	return &registerUsecase{
		customers: customers,
		admins:    admins,
		waiters:   waiters,
		log:       log,
		secretKey: secretKey,
	}
}

func (u *registerUsecase) RegisterCustomer(ctx context.Context, dto *domain.RegisterCustomerDTO) (string, error) {
	const op = "register.Customer"
	log := u.log.With(slog.String("op", op), slog.String("email", dto.Email))

	log.Info("registering customer")

	isExists, err := u.customers.CheckEmailExists(ctx, dto.Email)
	if err != nil {
		log.Error("failed to check email exists", "error", err)
		return "", err
	}

	if isExists {
		log.Info("customer with this email already exists")
		return "", er.ErrCustomerAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(dto.Password, bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", "error", err)
		return "", err
	}
	dto.Password = hashedPassword

	id, err := u.customers.CreateCustomer(ctx, dto)
	if err != nil {
		log.Error("failed to create customer", "error", err)
		return "", err
	}

	log.Info("customer registered", "customerId", id)
	return id, nil
}

func (u *registerUsecase) RegisterWaiter(ctx context.Context, dto *domain.RegisterWaiterDTO, token string) (string, error) {
	return "nil", status.Errorf(codes.Unimplemented, "method RegisterWaiter not implemented")
}

func (u *registerUsecase) RegisterAdmin(ctx context.Context, dto *domain.RegisterAdminDTO, token string) (string, error) {
	return "nil", status.Errorf(codes.Unimplemented, "method RegisterAdmin not implemented")
}
