package usecase

import (
	"context"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain"
)

type CustomerRepo interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateCustomer(ctx context.Context, dto *domain.RegisterCustomerDTO) (string, error)
}

type AdminRepo interface{}
type WaiterRepo interface{}
type TokensRepo interface{}
