package usecase_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/dto"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
}

func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type mockReservationRepo struct {
	mock.Mock
}

func (m *mockReservationRepo) CreateReservation(ctx context.Context, dto *dto.CreateReservationDTO) (uuid.UUID, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockReservationRepo) SetReservationStatus(ctx context.Context, reservationID uuid.UUID, status string) error {
	args := m.Called(ctx, reservationID, status)
	return args.Error(0)
}

func (m *mockReservationRepo) CloseEndedReservations(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockReservationRepo) GetTableExists(ctx context.Context, tableID uuid.UUID) (bool, error) {
	args := m.Called(ctx, tableID)
	return args.Get(0).(bool), args.Error(1)
}
