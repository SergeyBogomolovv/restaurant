package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/SergeyBogomolovv/restaurant/common/constants"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/errors"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReservationUsecase_CreateReservation(t *testing.T) {
	ctx := context.Background()
	logger := NewTestLogger()
	mockRepo := new(mockReservationRepo)
	mockBroker := new(mockBroker)

	usecase := usecase.NewReservationUsecase(logger, mockRepo, mockBroker)

	t.Run("success", func(t *testing.T) {
		tableId := uuid.New()
		resultId := uuid.New()

		mockRepo.On("GetTableExists", ctx, tableId).Return(true, nil).Once()
		mockRepo.On("CreateReservation", ctx, mock.Anything).Return(&dto.ReservationCreated{ReservationID: resultId}, nil).Once()
		mockBroker.On("Publish", "reservation.created", mock.Anything).Return(nil).Once()
		reservationID, err := usecase.CreateReservation(ctx, &dto.CreateReservationDTO{TableID: tableId})

		assert.NoError(t, err)
		assert.Equal(t, reservationID, resultId)
		mockRepo.AssertExpectations(t)
		mockBroker.AssertExpectations(t)
	})

	t.Run("table not found", func(t *testing.T) {
		tableId := uuid.New()

		mockRepo.On("GetTableExists", ctx, tableId).Return(false, nil).Once()

		id, err := usecase.CreateReservation(ctx, &dto.CreateReservationDTO{TableID: tableId})

		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, errs.ErrTableNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("table is reserved", func(t *testing.T) {
		tableId := uuid.New()
		mockRepo.On("GetTableExists", ctx, tableId).Return(true, nil).Once()
		mockRepo.On("CreateReservation", ctx, mock.Anything).Return((*dto.ReservationCreated)(nil), errs.ErrTableAlreadyReserved).Once()
		id, err := usecase.CreateReservation(ctx, &dto.CreateReservationDTO{TableID: tableId})

		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, errs.ErrTableAlreadyReserved)
		mockRepo.AssertExpectations(t)
	})
}

func TestReservationUsecase_CancelReservation(t *testing.T) {
	ctx := context.Background()
	logger := NewTestLogger()
	mockRepo := new(mockReservationRepo)
	mockBroker := new(mockBroker)

	usecase := usecase.NewReservationUsecase(logger, mockRepo, mockBroker)

	t.Run("succes", func(t *testing.T) {
		id := uuid.New()
		mockRepo.On("SetReservationStatus", ctx, id, constants.ReservationStatusCancelled).Return(nil).Once()

		err := usecase.CancelReservation(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("reservation not found", func(t *testing.T) {
		id := uuid.New()
		mockRepo.On("SetReservationStatus", ctx, id, constants.ReservationStatusCancelled).Return(errs.ErrReservationNotFound).Once()

		err := usecase.CancelReservation(ctx, id)

		assert.ErrorIs(t, err, errs.ErrReservationNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestReservationUsecase_CloseReservation(t *testing.T) {
	ctx := context.Background()
	logger := NewTestLogger()
	mockRepo := new(mockReservationRepo)
	mockBroker := new(mockBroker)

	usecase := usecase.NewReservationUsecase(logger, mockRepo, mockBroker)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		mockRepo.On("SetReservationStatus", ctx, id, constants.ReservationStatusClosed).Return(nil).Once()
		err := usecase.CloseReservation(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("reservation not found", func(t *testing.T) {
		id := uuid.New()
		mockRepo.On("SetReservationStatus", ctx, id, constants.ReservationStatusClosed).Return(errs.ErrReservationNotFound).Once()
		err := usecase.CloseReservation(ctx, id)

		assert.ErrorIs(t, err, errs.ErrReservationNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestReservationUsecase_CheckEndedReservations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := NewTestLogger()
	mockRepo := new(mockReservationRepo)

	tickerDuration := 50 * time.Millisecond
	usecase := usecase.NewReservationUsecase(logger, mockRepo, nil)

	mockRepo.On("CloseEndedReservations", mock.Anything).Return(int64(2), nil)

	go usecase.RunEndedReservationsChecker(ctx, tickerDuration)

	done := make(chan struct{})
	go func() {
		time.Sleep(2 * tickerDuration)
		cancel()
		close(done)
	}()
	<-done
	mockRepo.AssertCalled(t, "CloseEndedReservations", mock.Anything)
}
