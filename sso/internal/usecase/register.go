package usecase

import (
	"context"
	"log/slog"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Broker interface {
	Publish(routingKey string, payload any) error
}

type CustomerRegisterRepo interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateCustomer(ctx context.Context, dto *dto.CreateCustomerDTO) (*dto.RegisterCustomerResult, error)
}

type AdminRegisterRepo interface {
	CheckLoginExists(ctx context.Context, login string) (bool, error)
	CreateAdmin(ctx context.Context, dto *dto.CreateAdminDTO) (*dto.RegisterAdminResult, error)
}

type WaiterRegisterRepo interface {
	CheckLoginExists(ctx context.Context, login string) (bool, error)
	CreateWaiter(ctx context.Context, dto *dto.CreateWaiterDTO) (*dto.RegisterWaiterResult, error)
}

type registerUsecase struct {
	customers CustomerRegisterRepo
	admins    AdminRegisterRepo
	waiters   WaiterRegisterRepo
	log       *slog.Logger
	broker    Broker
	secretKey string
}

func NewRegisterUsecase(
	log *slog.Logger,
	customers CustomerRegisterRepo,
	waiters WaiterRegisterRepo,
	admins AdminRegisterRepo,
	broker Broker,
	secretKey string,
) *registerUsecase {
	return &registerUsecase{
		customers: customers,
		admins:    admins,
		waiters:   waiters,
		log:       log,
		secretKey: secretKey,
		broker:    broker,
	}
}

func (u *registerUsecase) RegisterCustomer(ctx context.Context, payload *dto.RegisterCustomerDTO) (uuid.UUID, error) {
	const op = "register.Customer"
	log := u.log.With(slog.String("op", op), slog.String("email", payload.Email))

	log.Info("registering customer")

	isExists, err := u.customers.CheckEmailExists(ctx, payload.Email)
	if err != nil {
		log.Error("failed to check email exists", "error", err)
		return uuid.Nil, err
	}

	if isExists {
		log.Info("customer with this email already exists")
		return uuid.Nil, errs.ErrCustomerAlreadyExists
	}

	hashedPassword, err := u.HashPassword(payload.Password)
	if err != nil {
		log.Error("failed to hash password", "error", err)
		return uuid.Nil, err
	}

	result, err := u.customers.CreateCustomer(ctx, &dto.CreateCustomerDTO{
		Email:     payload.Email,
		Name:      payload.Name,
		Birthdate: payload.Birthdate,
		Password:  hashedPassword,
	})
	if err != nil {
		log.Error("failed to create customer", "error", err)
		return uuid.Nil, err
	}

	if err := u.broker.Publish("register.customer", result); err != nil {
		log.Error("failed to publish message", "error", err)
		return uuid.Nil, err
	}

	log.Info("customer registered", "customerId", result.CustomerID)
	return result.CustomerID, nil
}

func (u *registerUsecase) RegisterWaiter(ctx context.Context, payload *dto.RegisterWaiterDTO) (uuid.UUID, error) {
	const op = "register.Waiter"
	log := u.log.With(slog.String("op", op), slog.String("login", payload.Login))

	log.Info("registering waiter")

	if payload.Token != u.secretKey {
		log.Info("invalid secret token")
		return uuid.Nil, errs.ErrInvalidSecretToken
	}
	isExists, err := u.waiters.CheckLoginExists(ctx, payload.Login)
	if err != nil {
		log.Error("failed to check login exists", "error", err)
		return uuid.Nil, err
	}

	if isExists {
		log.Info("waiter with this login already exists")
		return uuid.Nil, errs.ErrWaiterAlreadyExists
	}

	hashedPassword, err := u.HashPassword(payload.Password)
	if err != nil {
		log.Error("failed to hash password", "error", err)
		return uuid.Nil, err
	}

	result, err := u.waiters.CreateWaiter(ctx, &dto.CreateWaiterDTO{
		Login:     payload.Login,
		Password:  hashedPassword,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	})
	if err != nil {
		log.Error("failed to create waiter", "error", err)
		return uuid.Nil, err
	}

	if err := u.broker.Publish("register.waiter", result); err != nil {
		log.Error("failed to publish message", "error", err)
		return uuid.Nil, err
	}

	log.Info("waiter registered", "waiterId", result.WaiterID)
	return result.WaiterID, nil
}

func (u *registerUsecase) RegisterAdmin(ctx context.Context, payload *dto.RegisterAdminDTO) (uuid.UUID, error) {
	const op = "register.Admin"
	log := u.log.With(slog.String("op", op), slog.String("login", payload.Login))

	log.Info("registering admin")

	if payload.Token != u.secretKey {
		log.Info("invalid secret token")
		return uuid.Nil, errs.ErrInvalidSecretToken
	}

	isExists, err := u.admins.CheckLoginExists(ctx, payload.Login)
	if err != nil {
		log.Error("failed to check login exists", "error", err)
		return uuid.Nil, err
	}

	if isExists {
		log.Info("admin with this login already exists")
		return uuid.Nil, errs.ErrAdminAlreadyExists
	}

	hashedPassword, err := u.HashPassword(payload.Password)
	if err != nil {
		log.Error("failed to hash password", "error", err)
		return uuid.Nil, err
	}

	result, err := u.admins.CreateAdmin(ctx, &dto.CreateAdminDTO{
		Login:    payload.Login,
		Password: hashedPassword,
		Note:     payload.Note,
	})
	if err != nil {
		log.Error("failed to create admin", "error", err)
		return uuid.Nil, err
	}

	if err := u.broker.Publish("register.admin", result); err != nil {
		log.Error("failed to publish message", "error", err)
		return uuid.Nil, err
	}

	log.Info("admin registered", "adminId", result.AdminID)
	return result.AdminID, nil
}

func (u *registerUsecase) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
