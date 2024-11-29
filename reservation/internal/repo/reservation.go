package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/dto"
	errs "github.com/SergeyBogomolovv/restaurant/reservation/internal/domain/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type reservationRepo struct {
	db *sqlx.DB
}

func NewReservationRepo(db *sqlx.DB) *reservationRepo {
	return &reservationRepo{db: db}
}

func (r *reservationRepo) GetTableExists(ctx context.Context, tableID uuid.UUID) (bool, error) {
	var isExists bool
	if err := r.db.GetContext(ctx, &isExists, "SELECT TRUE FROM tables WHERE table_id = $1", tableID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isExists, nil
}

func (r *reservationRepo) CreateReservation(ctx context.Context, dto *dto.CreateReservationDTO) (uuid.UUID, error) {
	var id uuid.UUID
	query := `INSERT INTO reservations (customer_id, table_id, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING reservation_id`
	if err := r.db.GetContext(ctx, &id, query, dto.CustomerID, dto.TableID, dto.StartTime, dto.EndTime); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Message == "table already reserved" {
				return uuid.Nil, errs.ErrTableAlreadyReserved
			}
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (r *reservationRepo) SetReservationStatus(ctx context.Context, reservationID uuid.UUID, status string) error {
	query := `UPDATE reservations SET status = $1 WHERE reservation_id = $2`
	res, err := r.db.ExecContext(ctx, query, status, reservationID)
	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return errs.ErrReservationNotFound
	}
	return nil
}

func (r *reservationRepo) CloseEndedReservations(ctx context.Context) (int64, error) {
	query := `UPDATE reservations SET status = 'closed' WHERE end_time < now() AND status = 'active'`
	res, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
