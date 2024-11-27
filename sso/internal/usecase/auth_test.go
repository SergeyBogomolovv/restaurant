package usecase_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/SergeyBogomolovv/restaurant/common/constants"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/entities"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/usecase"
	"github.com/SergeyBogomolovv/restaurant/sso/pkg/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type mockCustomerAuthRepo struct {
	mock.Mock
}

func (m *mockCustomerAuthRepo) GetCustomerByEmail(ctx context.Context, email string) (*entities.CustomerEntity, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CustomerEntity), args.Error(1)
}

type mockAdminAuthRepo struct {
	mock.Mock
}

func (m *mockAdminAuthRepo) GetAdminByLogin(ctx context.Context, login string) (*entities.AdminEntity, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AdminEntity), args.Error(1)
}

type mockWaiterAuthRepo struct {
	mock.Mock
}

func (m *mockWaiterAuthRepo) GetWaiterByLogin(ctx context.Context, login string) (*entities.WaiterEntity, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.WaiterEntity), args.Error(1)
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

func TestAuthUsecase_LoginCustomer(t *testing.T) {
	ctx := context.Background()
	mockCustomerRepo := new(mockCustomerAuthRepo)
	mockTokensRepo := new(mockTokensRepo)

	usecase := usecase.NewAuthUsecase(NewTestLogger(), mockCustomerRepo, nil, nil, mockTokensRepo)

	t.Run("success", func(t *testing.T) {
		password := "password123"
		email := "test@example.com"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		customer := &entities.CustomerEntity{
			CustomerID: "123",
			Password:   hashedPassword,
		}

		mockCustomerRepo.On("GetCustomerByEmail", ctx, email).Return(customer, nil)
		mockTokensRepo.On("SignAccessToken", customer.CustomerID, constants.RoleCustomer).Return("access-token", nil)
		mockTokensRepo.On("GenerateRefreshToken", ctx, customer.CustomerID, constants.RoleCustomer).Return("refresh-token", nil)

		loginDTO := &dto.LoginCustomerDTO{Email: email, Password: password}

		tokens, err := usecase.LoginCustomer(ctx, loginDTO)

		assert.NoError(t, err)
		assert.Equal(t, "access-token", tokens.AccessToken)
		assert.Equal(t, "refresh-token", tokens.RefreshToken)
		mockCustomerRepo.AssertExpectations(t)
		mockTokensRepo.AssertExpectations(t)
	})

	t.Run("customer not found", func(t *testing.T) {
		email := "notfound@example.com"
		mockCustomerRepo.On("GetCustomerByEmail", ctx, email).Return(nil, errs.ErrCustomerNotFound)

		loginDTO := &dto.LoginCustomerDTO{Email: email, Password: "password123"}
		tokens, err := usecase.LoginCustomer(ctx, loginDTO)

		assert.Nil(t, tokens)
		assert.ErrorIs(t, err, errs.ErrInvalidCredentials)
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		email := "wrongpassword@example.com"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		customer := &entities.CustomerEntity{
			CustomerID: "123",
			Password:   hashedPassword,
		}
		mockCustomerRepo.On("GetCustomerByEmail", ctx, email).Return(customer, nil)

		loginDTO := &dto.LoginCustomerDTO{Email: email, Password: "wrongpassword"}
		tokens, err := usecase.LoginCustomer(ctx, loginDTO)

		assert.Nil(t, tokens)
		assert.ErrorIs(t, err, errs.ErrInvalidCredentials)
		mockCustomerRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_LoginAdmin(t *testing.T) {
	ctx := context.Background()
	mockAdminRepo := new(mockAdminAuthRepo)
	mockTokensRepo := new(mockTokensRepo)

	usecase := usecase.NewAuthUsecase(NewTestLogger(), nil, nil, mockAdminRepo, mockTokensRepo)

	t.Run("success", func(t *testing.T) {
		password := "adminpass123"
		login := "adminlogin"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		admin := &entities.AdminEntity{
			AdminID:  "456",
			Password: hashedPassword,
		}

		mockAdminRepo.On("GetAdminByLogin", ctx, login).Return(admin, nil)
		mockTokensRepo.On("SignAccessToken", admin.AdminID, constants.RoleAdmin).Return("admin-access-token", nil)
		mockTokensRepo.On("GenerateRefreshToken", ctx, admin.AdminID, constants.RoleAdmin).Return("admin-refresh-token", nil)

		loginDTO := &dto.LoginEmployeeDTO{Login: login, Password: password}

		tokens, err := usecase.LoginAdmin(ctx, loginDTO)

		assert.NoError(t, err)
		assert.Equal(t, "admin-access-token", tokens.AccessToken)
		assert.Equal(t, "admin-refresh-token", tokens.RefreshToken)
		mockAdminRepo.AssertExpectations(t)
		mockTokensRepo.AssertExpectations(t)
	})

	t.Run("admin not found", func(t *testing.T) {
		login := "unknownadmin"
		mockAdminRepo.On("GetAdminByLogin", ctx, login).Return(nil, errs.ErrAdminNotFound)

		loginDTO := &dto.LoginEmployeeDTO{Login: login, Password: "wrongpass"}
		tokens, err := usecase.LoginAdmin(ctx, loginDTO)

		assert.Nil(t, tokens)
		assert.ErrorIs(t, err, errs.ErrInvalidCredentials)
		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		login := "adminlogin"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
		admin := &entities.AdminEntity{
			AdminID:  "456",
			Password: hashedPassword,
		}
		mockAdminRepo.On("GetAdminByLogin", ctx, login).Return(admin, nil)

		loginDTO := &dto.LoginEmployeeDTO{Login: login, Password: "wrongpass"}
		tokens, err := usecase.LoginAdmin(ctx, loginDTO)

		assert.Nil(t, tokens)
		assert.ErrorIs(t, err, errs.ErrInvalidCredentials)
		mockAdminRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_LoginWaiter(t *testing.T) {
	ctx := context.Background()
	mockWaiterRepo := new(mockWaiterAuthRepo)
	mockTokensRepo := new(mockTokensRepo)

	usecase := usecase.NewAuthUsecase(NewTestLogger(), nil, mockWaiterRepo, nil, mockTokensRepo)

	t.Run("success", func(t *testing.T) {
		password := "waiterpass123"
		login := "waiterlogin"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		waiter := &entities.WaiterEntity{
			WaiterID: "789",
			Password: hashedPassword,
		}

		mockWaiterRepo.On("GetWaiterByLogin", ctx, login).Return(waiter, nil)
		mockTokensRepo.On("SignAccessToken", waiter.WaiterID, constants.RoleWaiter).Return("waiter-access-token", nil)
		mockTokensRepo.On("GenerateRefreshToken", ctx, waiter.WaiterID, constants.RoleWaiter).Return("waiter-refresh-token", nil)

		loginDTO := &dto.LoginEmployeeDTO{Login: login, Password: password}

		tokens, err := usecase.LoginWaiter(ctx, loginDTO)

		assert.NoError(t, err)
		assert.Equal(t, "waiter-access-token", tokens.AccessToken)
		assert.Equal(t, "waiter-refresh-token", tokens.RefreshToken)
		mockWaiterRepo.AssertExpectations(t)
		mockTokensRepo.AssertExpectations(t)
	})

	t.Run("waiter not found", func(t *testing.T) {
		login := "unknownwaiter"
		mockWaiterRepo.On("GetWaiterByLogin", ctx, login).Return(nil, errs.ErrWaiterNotFound)

		loginDTO := &dto.LoginEmployeeDTO{Login: login, Password: "password"}
		tokens, err := usecase.LoginWaiter(ctx, loginDTO)

		assert.Nil(t, tokens)
		assert.ErrorIs(t, err, errs.ErrInvalidCredentials)
		mockWaiterRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		login := "waiterlogin"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
		waiter := &entities.WaiterEntity{
			WaiterID: "789",
			Password: hashedPassword,
		}
		mockWaiterRepo.On("GetWaiterByLogin", ctx, login).Return(waiter, nil)

		loginDTO := &dto.LoginEmployeeDTO{Login: login, Password: "wrongpass"}
		tokens, err := usecase.LoginWaiter(ctx, loginDTO)

		assert.Nil(t, tokens)
		assert.ErrorIs(t, err, errs.ErrInvalidCredentials)
		mockWaiterRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_Refresh(t *testing.T) {
	ctx := context.Background()
	mockTokensRepo := new(mockTokensRepo)

	usecase := usecase.NewAuthUsecase(NewTestLogger(), nil, nil, nil, mockTokensRepo)

	t.Run("success", func(t *testing.T) {
		refreshPayload := &payload.JwtPayload{
			EntityID: "123",
			Role:     constants.RoleCustomer,
		}

		mockTokensRepo.On("VerifyRefreshToken", ctx, "valid-token").Return(refreshPayload, nil)
		mockTokensRepo.On("SignAccessToken", refreshPayload.EntityID, constants.RoleCustomer).Return("new-access-token", nil)

		newAccessToken, err := usecase.Refresh(ctx, "valid-token")

		assert.NoError(t, err)
		assert.Equal(t, "new-access-token", newAccessToken)
		mockTokensRepo.AssertExpectations(t)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		mockTokensRepo.On("VerifyRefreshToken", ctx, "invalid-token").Return(nil, errs.ErrInvalidJwtToken)

		newAccessToken, err := usecase.Refresh(ctx, "invalid-token")

		assert.Empty(t, newAccessToken)
		assert.ErrorIs(t, err, errs.ErrInvalidJwtToken)
		mockTokensRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_Logout(t *testing.T) {
	ctx := context.Background()
	mockTokensRepo := new(mockTokensRepo)

	mockTokensRepo.On("RevokeRefreshToken", ctx, "valid-token").Return(nil)

	usecase := usecase.NewAuthUsecase(NewTestLogger(), nil, nil, nil, mockTokensRepo)

	err := usecase.Logout(ctx, "valid-token")

	assert.NoError(t, err)
	mockTokensRepo.AssertExpectations(t)
}
