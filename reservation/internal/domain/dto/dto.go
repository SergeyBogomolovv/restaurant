package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateReservationDTO struct {
	CustomerID uuid.UUID `validate:"required,uuid"`
	TableID    uuid.UUID `validate:"required,uuid"`
	StartTime  time.Time `validate:"required"`
	EndTime    time.Time `validate:"required"`
}
