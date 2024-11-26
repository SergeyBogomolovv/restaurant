package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type waiterRepo struct {
	db *sqlx.DB
}

func NewWaiterRepo(db *sqlx.DB) *waiterRepo {
	return &waiterRepo{db: db}
}

func (r *waiterRepo) CheckLoginExists(ctx context.Context, login string) (bool, error) {
	var isExists bool
	if err := r.db.GetContext(ctx, &isExists, "SELECT TRUE FROM waiters WHERE login = $1", login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isExists, nil
}

func (r *waiterRepo) CreateWaiter(ctx context.Context, dto *dto.CreateWaiterDTO) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.db.GetContext(ctx, &id, `
		INSERT INTO waiters (login, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING waiter_id
		`, dto.Login, dto.Password, dto.FirstName, dto.LastName); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
