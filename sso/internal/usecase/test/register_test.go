package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/entities"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUsecase_RegisterCustomer(t *testing.T) {
	ctx := context.Background()
	log := NewTestLogger()
	customerRepo := new(mockCustomerRegisterRepo)
	usecase := usecase.NewRegisterUsecase(log, customerRepo, nil, nil, "secretKey")

	payload := &dto.RegisterCustomerDTO{
		Email:     "test@example.com",
		Password:  "password123",
		Name:      "John Doe",
		Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			customerRepo.ExpectedCalls = nil
			customerRepo.Calls = nil
		})
		customerRepo.On("CheckEmailExists", ctx, payload.Email).Return(false, nil)
		customerId := uuid.New()
		customerRepo.On("CreateCustomer", ctx, mock.Anything).Return(&entities.Customer{CustomerID: customerId}, nil)

		id, err := usecase.RegisterCustomer(ctx, payload)
		assert.NoError(t, err)
		assert.Equal(t, customerId, id)

		customerRepo.AssertExpectations(t)
	})

	t.Run("email exists", func(t *testing.T) {
		t.Cleanup(func() {
			customerRepo.ExpectedCalls = nil
			customerRepo.Calls = nil
		})
		customerRepo.On("CheckEmailExists", ctx, payload.Email).Return(true, nil)

		id, err := usecase.RegisterCustomer(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrCustomerAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		customerRepo.AssertExpectations(t)
	})

}

func TestRegisterUsecase_RegisterWaiter(t *testing.T) {
	ctx := context.Background()
	log := NewTestLogger()
	waiterRepo := new(mockWaiterRegisterRepo)
	usecase := usecase.NewRegisterUsecase(log, nil, waiterRepo, nil, "secretKey")

	payload := &dto.RegisterWaiterDTO{
		Login:     "waiter123",
		Password:  "password123",
		FirstName: "Jane",
		LastName:  "Doe",
		Token:     "secretKey",
	}

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			waiterRepo.ExpectedCalls = nil
			waiterRepo.Calls = nil
		})
		waiterRepo.On("CheckLoginExists", ctx, payload.Login).Return(false, nil)
		waiterID := uuid.New()
		waiterRepo.On("CreateWaiter", ctx, mock.Anything).Return(&entities.Waiter{WaiterID: waiterID}, nil)

		id, err := usecase.RegisterWaiter(ctx, payload)
		assert.NoError(t, err)
		assert.Equal(t, id, waiterID)

		waiterRepo.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Cleanup(func() {
			waiterRepo.ExpectedCalls = nil
			waiterRepo.Calls = nil
		})
		payload := &dto.RegisterWaiterDTO{
			Token: "invalidToken",
		}

		id, err := usecase.RegisterWaiter(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrInvalidSecretToken)
		assert.Equal(t, id, uuid.Nil)
	})

	t.Run("login exists", func(t *testing.T) {
		t.Cleanup(func() {
			waiterRepo.ExpectedCalls = nil
			waiterRepo.Calls = nil
		})
		waiterRepo.On("CheckLoginExists", ctx, payload.Login).Return(true, nil)

		id, err := usecase.RegisterWaiter(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrWaiterAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		waiterRepo.AssertExpectations(t)
	})
}

func TestRegisterUsecase_RegisterAdmin(t *testing.T) {
	ctx := context.Background()
	log := NewTestLogger()
	adminRepo := new(mockAdminRegisterRepo)
	usecase := usecase.NewRegisterUsecase(log, nil, nil, adminRepo, "secretKey")

	payload := &dto.RegisterAdminDTO{
		Login:    "admin123",
		Password: "password123",
		Note:     "Super admin",
		Token:    "secretKey",
	}

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			adminRepo.ExpectedCalls = nil
			adminRepo.Calls = nil
		})
		adminRepo.On("CheckLoginExists", ctx, payload.Login).Return(false, nil)
		adminId := uuid.New()
		adminRepo.On("CreateAdmin", ctx, mock.Anything).Return(&entities.Admin{AdminID: adminId}, nil)

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.NoError(t, err)
		assert.Equal(t, id, adminId)

		adminRepo.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Cleanup(func() {
			adminRepo.ExpectedCalls = nil
			adminRepo.Calls = nil
		})
		payload := &dto.RegisterAdminDTO{
			Token: "invalidToken",
		}

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrInvalidSecretToken)
		assert.Equal(t, id, uuid.Nil)
		adminRepo.AssertExpectations(t)
	})

	t.Run("login exists", func(t *testing.T) {
		t.Cleanup(func() {
			adminRepo.ExpectedCalls = nil
			adminRepo.Calls = nil
		})
		adminRepo.On("CheckLoginExists", ctx, payload.Login).Return(true, nil)

		id, err := usecase.RegisterAdmin(ctx, payload)
		assert.ErrorIs(t, err, errs.ErrAdminAlreadyExists)
		assert.Equal(t, id, uuid.Nil)

		adminRepo.AssertExpectations(t)
	})
}
