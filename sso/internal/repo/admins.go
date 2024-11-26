package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type adminRepo struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) CheckLoginExists(ctx context.Context, login string) (bool, error) {
	var isExists bool
	if err := r.db.GetContext(ctx, &isExists, "SELECT TRUE FROM admins WHERE login = $1", login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isExists, nil
}

func (r *adminRepo) CreateAdmin(ctx context.Context, dto *dto.CreateAdminDTO) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.db.GetContext(ctx, &id, `
		INSERT INTO admins (login, password, note) VALUES ($1, $2, $3) RETURNING admin_id
		`, dto.Login, dto.Password, dto.Note); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
