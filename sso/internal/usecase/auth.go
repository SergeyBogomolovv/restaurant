package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/restaurant/common/constants"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/entities"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/SergeyBogomolovv/restaurant/sso/pkg/payload"
	"golang.org/x/crypto/bcrypt"
)

type CustomerAuthRepo interface {
	GetCustomerByEmail(ctx context.Context, email string) (*entities.CustomerEntity, error)
}

type AdminAuthRepo interface {
	GetAdminByLogin(ctx context.Context, login string) (*entities.AdminEntity, error)
}

type WaiterAuthRepo interface {
	GetWaiterByLogin(ctx context.Context, login string) (*entities.WaiterEntity, error)
}

type TokensRepo interface {
	GenerateRefreshToken(ctx context.Context, entityID string, role string) (string, error)
	VerifyRefreshToken(ctx context.Context, token string) (*payload.JwtPayload, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	SignAccessToken(entityID string, role string) (string, error)
}

type authUsecase struct {
	customers CustomerAuthRepo
	admins    AdminAuthRepo
	waiters   WaiterAuthRepo
	tokens    TokensRepo
	log       *slog.Logger
}

func NewAuthUsecase(log *slog.Logger, customers CustomerAuthRepo, waiters WaiterAuthRepo, admins AdminAuthRepo, tokens TokensRepo) *authUsecase {
	return &authUsecase{
		customers: customers,
		tokens:    tokens,
		admins:    admins,
		waiters:   waiters,
		log:       log,
	}
}

func (u *authUsecase) LoginCustomer(ctx context.Context, loginCustomerDto *dto.LoginCustomerDTO) (*dto.TokensDTO, error) {
	const op = "login.Customer"
	log := u.log.With(slog.String("op", op), slog.String("email", loginCustomerDto.Email))

	log.Info("logging in customer")

	customer, err := u.customers.GetCustomerByEmail(ctx, loginCustomerDto.Email)
	if err != nil {
		if errors.Is(err, errs.ErrCustomerNotFound) {
			log.Info("customer not found")
			return nil, errs.ErrInvalidCredentials
		}
		log.Error("failed to get customer by email", "error", err)
		return nil, err
	}

	if err := ComparePassword(customer.Password, loginCustomerDto.Password); err != nil {
		log.Info("invalid password")
		return nil, errs.ErrInvalidCredentials
	}

	accessToken, err := u.tokens.SignAccessToken(customer.CustomerID, constants.RoleCustomer)
	if err != nil {
		log.Error("failed to sign access token", "error", err)
		return nil, err
	}

	refreshToken, err := u.tokens.GenerateRefreshToken(ctx, customer.CustomerID, constants.RoleCustomer)
	if err != nil {
		log.Error("failed to generate refresh token", "error", err)
		return nil, err
	}

	return &dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (u *authUsecase) LoginAdmin(ctx context.Context, payload *dto.LoginEmployeeDTO) (*dto.TokensDTO, error) {
	const op = "login.Admin"
	log := u.log.With(slog.String("op", op), slog.String("login", payload.Login))

	log.Info("logging in admin")

	admin, err := u.admins.GetAdminByLogin(ctx, payload.Login)
	if err != nil {
		if errors.Is(err, errs.ErrAdminNotFound) {
			log.Info("admin not found")
			return nil, errs.ErrInvalidCredentials
		}
		log.Error("failed to get admin by login", "error", err)
		return nil, err
	}

	if err := ComparePassword(admin.Password, payload.Password); err != nil {
		log.Info("invalid password")
		return nil, errs.ErrInvalidCredentials
	}

	accessToken, err := u.tokens.SignAccessToken(admin.AdminID, constants.RoleAdmin)
	if err != nil {
		log.Error("failed to sign access token", "error", err)
		return nil, err
	}

	refreshToken, err := u.tokens.GenerateRefreshToken(ctx, admin.AdminID, constants.RoleAdmin)
	if err != nil {
		log.Error("failed to generate refresh token", "error", err)
		return nil, err
	}

	return &dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (u *authUsecase) LoginWaiter(ctx context.Context, payload *dto.LoginEmployeeDTO) (*dto.TokensDTO, error) {
	const op = "login.Waiter"
	log := u.log.With(slog.String("op", op), slog.String("login", payload.Login))

	log.Info("logging in waiter")

	waiter, err := u.waiters.GetWaiterByLogin(ctx, payload.Login)
	if err != nil {
		if errors.Is(err, errs.ErrWaiterNotFound) {
			log.Info("waiter not found")
			return nil, errs.ErrInvalidCredentials
		}
		log.Error("failed to get waiter by login", "error", err)
		return nil, err
	}

	if err := ComparePassword(waiter.Password, payload.Password); err != nil {
		log.Info("invalid password")
		return nil, errs.ErrInvalidCredentials
	}

	accessToken, err := u.tokens.SignAccessToken(waiter.WaiterID, constants.RoleWaiter)
	if err != nil {
		log.Error("failed to sign access token", "error", err)
		return nil, err
	}

	refreshToken, err := u.tokens.GenerateRefreshToken(ctx, waiter.WaiterID, constants.RoleWaiter)
	if err != nil {
		log.Error("failed to generate refresh token", "error", err)
		return nil, err
	}

	return &dto.TokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (u *authUsecase) Refresh(ctx context.Context, token string) (string, error) {
	const op = "auth.Refresh"
	log := u.log.With(slog.String("op", op))

	payload, err := u.tokens.VerifyRefreshToken(ctx, token)
	if err != nil {
		return "", errs.ErrInvalidJwtToken
	}

	accessToken, err := u.tokens.SignAccessToken(payload.EntityID, payload.Role)
	if err != nil {
		log.Error("failed to sign access token", "error", err)
		return "", err
	}

	return accessToken, nil
}

func (u *authUsecase) Logout(ctx context.Context, token string) error {
	const op = "auth.Logout"
	log := u.log.With(slog.String("op", op))

	if err := u.tokens.RevokeRefreshToken(ctx, token); err != nil {
		if errors.Is(err, errs.ErrInvalidJwtToken) {
			return nil
		}
		log.Error("failed to revoke refresh token", "error", err)
		return err
	}
	return nil
}

func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
