package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUsecase_RegisterCustomer(t *testing.T) {
	ctx := context.Background()
	customerRepo := new(mockCustomerRegisterRepo)
	usecase := usecase.NewRegisterUsecase(NewTestLogger(), customerRepo, nil, nil, nil, "secretKey")

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		result := &dto.RegisterCustomerResult{CustomerID: uuid.New()}

		customerRepo.On("CheckEmailExists", ctx, email).Return(false, nil).Once()
		customerRepo.On("CreateCustomerWithAction", ctx, mock.Anything, mock.Anything).Return(result, nil).Once()

		id, err := usecase.RegisterCustomer(ctx, &dto.RegisterCustomerDTO{Email: email})
		assert.NoError(t, err)
		assert.Equal(t, result.CustomerID, id)

		customerRepo.AssertExpectations(t)
	})

	t.Run("email exists", func(t *testing.T) {
		payload := &dto.RegisterCustomerDTO{Email: "test@example.com"}
		customerRepo.On("CheckEmailExists", ctx, payload.Email).Return(true, nil).Once()

		id, err := usecase.RegisterCustomer(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrCustomerAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		customerRepo.AssertExpectations(t)
	})

	t.Run("check email error", func(t *testing.T) {
		email := "error@example.com"
		simulatedErr := errors.New("database error")
		customerRepo.On("CheckEmailExists", ctx, email).Return(false, simulatedErr).Once()

		id, err := usecase.RegisterCustomer(ctx, &dto.RegisterCustomerDTO{Email: email})
		assert.Error(t, err)
		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, simulatedErr)

		customerRepo.AssertExpectations(t)
	})

	t.Run("create customer error", func(t *testing.T) {
		email := "fail@example.com"
		simulatedErr := errors.New("insert error")
		customerRepo.On("CheckEmailExists", ctx, email).Return(false, nil).Once()
		customerRepo.On("CreateCustomerWithAction",
			ctx,
			mock.Anything,
			mock.Anything).Return((*dto.RegisterCustomerResult)(nil), simulatedErr).Once()

		id, err := usecase.RegisterCustomer(ctx, &dto.RegisterCustomerDTO{Email: email})
		assert.Error(t, err)
		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, simulatedErr)

		customerRepo.AssertExpectations(t)
	})

	t.Run("broker publish error", func(t *testing.T) {
		email := "broker@example.com"
		simulatedErr := errors.New("broker error")

		customerRepo.On("CheckEmailExists", ctx, email).Return(false, nil).Once()
		customerRepo.On("CreateCustomerWithAction",
			ctx,
			mock.Anything,
			mock.Anything).Return((*dto.RegisterCustomerResult)(nil), simulatedErr).Once()

		id, err := usecase.RegisterCustomer(ctx, &dto.RegisterCustomerDTO{Email: email})
		assert.Error(t, err)
		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, simulatedErr)

		customerRepo.AssertExpectations(t)
	})
}

func TestRegisterUsecase_RegisterAdmin(t *testing.T) {
	ctx := context.Background()
	adminRepo := new(mockAdminRegisterRepo)
	key := "secretKey"
	usecase := usecase.NewRegisterUsecase(NewTestLogger(), nil, nil, adminRepo, nil, key)

	t.Run("success", func(t *testing.T) {
		result := &dto.RegisterAdminResult{AdminID: uuid.New()}
		payload := &dto.RegisterAdminDTO{Login: "admin123", Token: key}

		adminRepo.On("CheckLoginExists", ctx, payload.Login).Return(false, nil).Once()
		adminRepo.On("CreateAdminWithAction", ctx, mock.Anything, mock.Anything).
			Return(result, nil).
			Once()

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.NoError(t, err)
		assert.Equal(t, id, result.AdminID)

		adminRepo.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		payload := &dto.RegisterAdminDTO{Token: "invalidToken"}

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrInvalidSecretToken)
		assert.Equal(t, id, uuid.Nil)
	})

	t.Run("login exists", func(t *testing.T) {
		payload := &dto.RegisterAdminDTO{Login: "admin123", Token: key}
		adminRepo.On("CheckLoginExists", ctx, payload.Login).Return(true, nil).Once()

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrAdminAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		adminRepo.AssertExpectations(t)
	})

	t.Run("check login error", func(t *testing.T) {
		login := "error123"
		simulatedErr := errors.New("database error")
		adminRepo.On("CheckLoginExists", ctx, login).Return(false, simulatedErr).Once()

		id, err := usecase.RegisterAdmin(ctx, &dto.RegisterAdminDTO{Token: key, Login: login})
		assert.Error(t, err)
		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, simulatedErr)

		adminRepo.AssertExpectations(t)
	})
}

func TestRegisterUsecase_RegisterWaiter(t *testing.T) {
	ctx := context.Background()
	waiterRepo := new(mockWaiterRegisterRepo)
	key := "secretKey"
	usecase := usecase.NewRegisterUsecase(NewTestLogger(), nil, waiterRepo, nil, nil, key)

	t.Run("success", func(t *testing.T) {
		login := "waiter123"
		result := &dto.RegisterWaiterResult{WaiterID: uuid.New()}

		waiterRepo.On("CheckLoginExists", ctx, login).Return(false, nil).Once()
		waiterRepo.On("CreateWaiterWithAction", ctx, mock.Anything, mock.Anything).
			Return(result, nil).
			Once()

		id, err := usecase.RegisterWaiter(ctx, &dto.RegisterWaiterDTO{Token: key, Login: login})
		assert.NoError(t, err)
		assert.Equal(t, id, result.WaiterID)

		waiterRepo.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		payload := &dto.RegisterWaiterDTO{Token: "invalidToken"}

		id, err := usecase.RegisterWaiter(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrInvalidSecretToken)
		assert.Equal(t, id, uuid.Nil)
	})

	t.Run("login exists", func(t *testing.T) {
		payload := &dto.RegisterWaiterDTO{Login: "waiter123", Token: key}
		waiterRepo.On("CheckLoginExists", ctx, payload.Login).Return(true, nil).Once()

		id, err := usecase.RegisterWaiter(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrWaiterAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		waiterRepo.AssertExpectations(t)
	})

	t.Run("check login error", func(t *testing.T) {
		login := "error123"
		simulatedErr := errors.New("database error")
		waiterRepo.On("CheckLoginExists", ctx, login).Return(false, simulatedErr).Once()

		id, err := usecase.RegisterWaiter(ctx, &dto.RegisterWaiterDTO{Token: key, Login: login})
		assert.Error(t, err)
		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, simulatedErr)

		waiterRepo.AssertExpectations(t)
	})
}
