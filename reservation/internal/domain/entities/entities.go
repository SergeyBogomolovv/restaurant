package entities

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ReservationID uuid.UUID `db:"reservation_id"`
	CustomerID    uuid.UUID `db:"customer_id"`
	StartTime     time.Time `db:"start_time"`
	EndTime       time.Time `db:"end_time"`
	Status        string    `db:"status"`
	TableID       uuid.UUID `db:"table_id"`
	PersonsCount  int       `db:"persons_count"`
}

type Table struct {
	TableID     uuid.UUID `db:"table_id"`
	TableNumber int       `db:"table_number"`
	Capacity    int       `db:"capacity"`
}
