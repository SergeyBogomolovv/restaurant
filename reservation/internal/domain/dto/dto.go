package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateReservationDTO struct {
	CustomerID   uuid.UUID `validate:"required,uuid"`
	TableID      uuid.UUID `validate:"required,uuid"`
	StartTime    time.Time `validate:"required"`
	EndTime      time.Time `validate:"required"`
	PersonsCount int       `validate:"required,min=1,max=8"`
}

type ReservationCreated struct {
	ReservationID uuid.UUID `json:"reservation_id" db:"reservation_id"`
	TableID       uuid.UUID `json:"table_id" db:"table_id"`
	PersonsCount  int       `json:"persons_count" db:"persons_count"`
	StartTime     time.Time `json:"start_time" db:"start_time"`
	EndTime       time.Time `json:"end_time" db:"end_time"`
}
