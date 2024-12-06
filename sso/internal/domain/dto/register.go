package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterCustomerDTO struct {
	Email     string    `validate:"required,email"`
	Name      string    `validate:"required"`
	Birthdate time.Time `validate:"required"`
	Password  string    `validate:"required"`
}

type RegisterCustomerResult struct {
	CustomerID uuid.UUID `json:"customer_id" db:"customer_id"`
	Name       string    `json:"name" db:"name"`
	Birthdate  string    `json:"birth_date" db:"birth_date"`
}

type RegisterAdminDTO struct {
	Note     string
	Login    string `validate:"required"`
	Password string `validate:"required"`
	Token    string `validate:"required"`
}

type RegisterWaiterResult struct {
	WaiterID  uuid.UUID `json:"waiter_id" db:"waiter_id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
}

type RegisterWaiterDTO struct {
	Login     string `validate:"required"`
	Password  string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Token     string `validate:"required"`
}

type RegisterAdminResult struct {
	AdminID uuid.UUID `json:"admin_id" db:"admin_id"`
	Login   string    `json:"login" db:"login"`
}
