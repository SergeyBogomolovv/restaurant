package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain"
	"github.com/jmoiron/sqlx"
)

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *customerRepo {
	return &customerRepo{db: db}
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

func (r *customerRepo) CreateCustomer(ctx context.Context, dto *domain.RegisterCustomerDTO) (string, error) {
	var customerId string
	query := "INSERT INTO customers (email, password, name, birth_date) VALUES ($1, $2, $3, $4) RETURNING customer_id"
	if err := r.db.GetContext(ctx, &customerId, query, dto.Email, dto.Password, dto.Name, dto.Birthdate); err != nil {
		return "", err
	}
	return customerId, nil
}
