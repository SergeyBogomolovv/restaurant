package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/entities"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/jmoiron/sqlx"
)

type waiterRepo struct {
	db *sqlx.DB
}

func NewWaiterRepo(db *sqlx.DB) *waiterRepo {
	return &waiterRepo{db: db}
}

func (r *waiterRepo) GetWaiterByLogin(ctx context.Context, login string) (*entities.Waiter, error) {
	waiter := new(entities.Waiter)
	if err := r.db.GetContext(ctx, waiter, "SELECT * FROM waiters WHERE login = $1", login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrWaiterNotFound
		}
		return nil, err
	}
	return waiter, nil
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

func (r *waiterRepo) CreateWaiter(ctx context.Context, dto *dto.CreateWaiterDTO) (*entities.Waiter, error) {
	waiter := new(entities.Waiter)
	if err := r.db.GetContext(ctx, waiter, `
		INSERT INTO waiters (login, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING *
		`, dto.Login, dto.Password, dto.FirstName, dto.LastName); err != nil {
		return nil, err
	}
	return waiter, nil
}
