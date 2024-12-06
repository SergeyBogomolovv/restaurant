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

type adminRepo struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) GetAdminByLogin(ctx context.Context, login string) (*entities.Admin, error) {
	admin := new(entities.Admin)
	if err := r.db.GetContext(ctx, admin, "SELECT * FROM admins WHERE login = $1", login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrAdminNotFound
		}
		return nil, err
	}
	return admin, nil
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

func (r *adminRepo) CreateAdmin(ctx context.Context, payload *dto.CreateAdminDTO) (*dto.RegisterAdminResult, error) {
	result := new(dto.RegisterAdminResult)
	if err := r.db.GetContext(ctx, result, `
		INSERT INTO admins (login, password, note)
		VALUES ($1, $2, $3) 
		RETURNING admin_id, login
		`, payload.Login, payload.Password, payload.Note); err != nil {
		return nil, err
	}
	return result, nil
}
