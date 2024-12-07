package usecase_test

import (
	"context"
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
	mockBroker := new(mockBroker)
	usecase := usecase.NewRegisterUsecase(NewTestLogger(), customerRepo, nil, nil, mockBroker, "secretKey")

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		result := &dto.RegisterCustomerResult{CustomerID: uuid.New()}

		customerRepo.On("CheckEmailExists", ctx, email).Return(false, nil).Once()
		customerRepo.On("CreateCustomer", ctx, mock.Anything).Return(result, nil).Once()
		mockBroker.On("Publish", "register.customer", result).Return(nil).Once()

		id, err := usecase.RegisterCustomer(ctx, &dto.RegisterCustomerDTO{Email: email})
		assert.NoError(t, err)
		assert.Equal(t, result.CustomerID, id)

		customerRepo.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})

	t.Run("email exists", func(t *testing.T) {
		payload := &dto.RegisterCustomerDTO{Email: "test@example.com"}
		customerRepo.On("CheckEmailExists", ctx, payload.Email).Return(true, nil).Once()

		id, err := usecase.RegisterCustomer(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrCustomerAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		customerRepo.AssertExpectations(t)
	})

}

func TestRegisterUsecase_RegisterWaiter(t *testing.T) {
	ctx := context.Background()
	waiterRepo := new(mockWaiterRegisterRepo)
	mockBroker := new(mockBroker)
	key := "secretKey"
	usecase := usecase.NewRegisterUsecase(NewTestLogger(), nil, waiterRepo, nil, mockBroker, key)

	t.Run("success", func(t *testing.T) {
		login := "waiter123"
		result := &dto.RegisterWaiterResult{WaiterID: uuid.New()}

		waiterRepo.On("CheckLoginExists", ctx, login).Return(false, nil).Once()
		waiterRepo.On("CreateWaiter", ctx, mock.Anything).Return(result, nil).Once()
		mockBroker.On("Publish", "register.waiter", result).Return(nil).Once()

		id, err := usecase.RegisterWaiter(ctx, &dto.RegisterWaiterDTO{Token: key, Login: login})
		assert.NoError(t, err)
		assert.Equal(t, id, result.WaiterID)

		waiterRepo.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
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
}

func TestRegisterUsecase_RegisterAdmin(t *testing.T) {
	ctx := context.Background()
	adminRepo := new(mockAdminRegisterRepo)
	mockBroker := new(mockBroker)
	key := "secretKey"
	usecase := usecase.NewRegisterUsecase(NewTestLogger(), nil, nil, adminRepo, mockBroker, key)

	t.Run("success", func(t *testing.T) {
		result := &dto.RegisterAdminResult{AdminID: uuid.New()}
		payload := &dto.RegisterAdminDTO{Login: "admin123", Token: key}

		adminRepo.On("CheckLoginExists", ctx, payload.Login).Return(false, nil).Once()
		adminRepo.On("CreateAdmin", ctx, mock.Anything).Return(result, nil).Once()
		mockBroker.On("Publish", "register.admin", result).Return(nil).Once()

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.NoError(t, err)
		assert.Equal(t, id, result.AdminID)

		adminRepo.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		payload := &dto.RegisterAdminDTO{Token: "invalidToken"}

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrInvalidSecretToken)
		assert.Equal(t, id, uuid.Nil)
		adminRepo.AssertExpectations(t)
	})

	t.Run("login exists", func(t *testing.T) {
		payload := &dto.RegisterAdminDTO{Login: "admin123", Token: key}
		adminRepo.On("CheckLoginExists", ctx, payload.Login).Return(true, nil)

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrAdminAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		adminRepo.AssertExpectations(t)
	})
}
