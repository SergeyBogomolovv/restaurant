package dto

import (
	"time"

	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/entities"
	"github.com/google/uuid"
)

type CreateReservationDTO struct {
	CustomerID   uuid.UUID `validate:"required,uuid"`
	TableID      uuid.UUID `validate:"required,uuid"`
	StartTime    time.Time `validate:"required"`
	EndTime      time.Time `validate:"required"`
	PersonsCount int       `validate:"required,min=1,max=8"`
}

type ReservationCreatedDTO struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	TableID       uuid.UUID `json:"table_id"`
	PersonsCount  int       `json:"persons_count"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

func NewReservationCreatedDTO(reservation *entities.Reservation) *ReservationCreatedDTO {
	return &ReservationCreatedDTO{
		ReservationID: reservation.ReservationID,
		TableID:       reservation.TableID,
		PersonsCount:  reservation.PersonsCount,
		StartTime:     reservation.StartTime,
		EndTime:       reservation.EndTime,
	}
}
