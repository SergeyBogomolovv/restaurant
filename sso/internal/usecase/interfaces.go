package usecase

import (
	"context"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain"
)

type CustomerRepo interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateCustomer(ctx context.Context, dto *domain.RegisterCustomerDTO) (string, error)
}

type AdminRepo interface {
	CheckLoginExists(ctx context.Context, login string) (bool, error)
	CreateAdmin(ctx context.Context, dto *domain.RegisterAdminDTO) (string, error)
}

type WaiterRepo interface {
	CheckLoginExists(ctx context.Context, login string) (bool, error)
	CreateWaiter(ctx context.Context, dto *domain.RegisterWaiterDTO) (string, error)
}

type TokensRepo interface{}
