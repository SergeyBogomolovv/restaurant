package usecase_test

import (
	"context"
	"io"
	"log/slog"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/entities"
	"github.com/SergeyBogomolovv/restaurant/sso/pkg/payload"
	"github.com/stretchr/testify/mock"
)

func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type mockCustomerAuthRepo struct {
	mock.Mock
}

func (m *mockCustomerAuthRepo) GetCustomerByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Customer), args.Error(1)
}

type mockAdminAuthRepo struct {
	mock.Mock
}

func (m *mockAdminAuthRepo) GetAdminByLogin(ctx context.Context, login string) (*entities.Admin, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Admin), args.Error(1)
}

type mockWaiterAuthRepo struct {
	mock.Mock
}

func (m *mockWaiterAuthRepo) GetWaiterByLogin(ctx context.Context, login string) (*entities.Waiter, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Waiter), args.Error(1)
}

type mockTokensRepo struct {
	mock.Mock
}

func (m *mockTokensRepo) GenerateRefreshToken(ctx context.Context, entityID string, role string) (string, error) {
	args := m.Called(ctx, entityID, role)
	return args.String(0), args.Error(1)
}

func (m *mockTokensRepo) VerifyRefreshToken(ctx context.Context, token string) (*payload.JwtPayload, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payload.JwtPayload), args.Error(1)
}

func (m *mockTokensRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	return m.Called(ctx, token).Error(0)
}

func (m *mockTokensRepo) SignAccessToken(entityID string, role string) (string, error) {
	args := m.Called(entityID, role)
	return args.String(0), args.Error(1)
}

type mockCustomerRegisterRepo struct {
	mock.Mock
}

func (m *mockCustomerRegisterRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *mockCustomerRegisterRepo) CreateCustomer(ctx context.Context, dto *dto.CreateCustomerDTO) (*entities.Customer, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*entities.Customer), args.Error(1)
}

type mockAdminRegisterRepo struct {
	mock.Mock
}

func (m *mockAdminRegisterRepo) CheckLoginExists(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func (m *mockAdminRegisterRepo) CreateAdmin(ctx context.Context, dto *dto.CreateAdminDTO) (*entities.Admin, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*entities.Admin), args.Error(1)
}

type mockWaiterRegisterRepo struct {
	mock.Mock
}

func (m *mockWaiterRegisterRepo) CheckLoginExists(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func (m *mockWaiterRegisterRepo) CreateWaiter(ctx context.Context, dto *dto.CreateWaiterDTO) (*entities.Waiter, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*entities.Waiter), args.Error(1)
}
