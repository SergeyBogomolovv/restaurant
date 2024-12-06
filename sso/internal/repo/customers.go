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

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *customerRepo {
	return &customerRepo{db: db}
}

func (r *customerRepo) GetCustomerByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	customer := new(entities.Customer)
	if err := r.db.GetContext(ctx, customer, "SELECT * FROM customers WHERE email = $1", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCustomerNotFound
		}
		return nil, err
	}
	return customer, nil
}

func (r *customerRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var isExists bool
	if err := r.db.GetContext(ctx, &isExists, "SELECT TRUE FROM customers WHERE email = $1", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isExists, nil
}

func (r *customerRepo) CreateCustomer(ctx context.Context, payload *dto.CreateCustomerDTO) (*dto.RegisterCustomerResult, error) {
	result := new(dto.RegisterCustomerResult)
	if err := r.db.GetContext(ctx, result, `
	INSERT INTO customers (email, password, name, birth_date) 
	VALUES ($1, $2, $3, $4)
	RETURNING customer_id, name, birth_date
	`, payload.Email, payload.Password, payload.Name, payload.Birthdate); err != nil {
		return nil, err
	}
	return result, nil
}
