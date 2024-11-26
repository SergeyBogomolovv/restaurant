package usecase

import (
	"context"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/google/uuid"
)

type CustomerRepo interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateCustomer(ctx context.Context, dto *dto.CreateCustomerDTO) (uuid.UUID, error)
}

type AdminRepo interface {
	CheckLoginExists(ctx context.Context, login string) (bool, error)
	CreateAdmin(ctx context.Context, dto *dto.CreateAdminDTO) (uuid.UUID, error)
}

type WaiterRepo interface {
	CheckLoginExists(ctx context.Context, login string) (bool, error)
	CreateWaiter(ctx context.Context, dto *dto.CreateWaiterDTO) (uuid.UUID, error)
}

type TokensRepo interface{}
