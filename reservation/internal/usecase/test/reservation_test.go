package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/errors"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReservationUsecase_CreateReservation(t *testing.T) {
	ctx := context.Background()
	logger := NewTestLogger()
	mockRepo := new(mockReservationRepo)
	mockRepo.On("CloseEndedReservations", ctx).Return(int64(0), nil)

	usecase := usecase.NewReservationUsecase(logger, mockRepo, ctx)

	t.Run("success", func(t *testing.T) {
		tableId := uuid.New()
		dto := &dto.CreateReservationDTO{
			CustomerID: uuid.New(),
			TableID:    tableId,
			StartTime:  time.Unix(1730455200, 0),
			EndTime:    time.Unix(1730458800, 0),
		}
		mockRepo.On("GetTableExists", ctx, tableId).Return(true, nil)
		mockRepo.On("CreateReservation", ctx, dto).Return(uuid.New(), nil)

		reservationID, err := usecase.CreateReservation(ctx, dto)

		assert.NoError(t, err)
		assert.NotEqual(t, reservationID, uuid.Nil)
	})

	t.Run("table not found", func(t *testing.T) {
		tableId := uuid.New()
		dto := &dto.CreateReservationDTO{
			CustomerID: uuid.New(),
			TableID:    tableId,
			StartTime:  time.Unix(1730455200, 0),
			EndTime:    time.Unix(1730458800, 0),
		}

		mockRepo.On("GetTableExists", ctx, tableId).Return(false, nil)

		id, err := usecase.CreateReservation(ctx, dto)

		assert.Equal(t, id, uuid.Nil)
		assert.ErrorIs(t, err, errs.ErrTableNotFound)
	})
}
