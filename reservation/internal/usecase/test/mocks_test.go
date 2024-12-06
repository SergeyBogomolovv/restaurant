package usecase_test

import (
	"context"
	"io"
	"log/slog"

	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type mockReservationRepo struct {
	mock.Mock
}

func (m *mockReservationRepo) CreateReservation(ctx context.Context, payload *dto.CreateReservationDTO) (*dto.ReservationCreated, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(*dto.ReservationCreated), args.Error(1)
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

type mockBroker struct {
	mock.Mock
}

func (m *mockBroker) Publish(key string, payload any) error {
	args := m.Called(key, payload)
	return args.Error(0)
}
