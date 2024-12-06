package dto

import "time"

type RegisterCustomerDTO struct {
	Email     string    `validate:"required,email"`
	Name      string    `validate:"required"`
	Birthdate time.Time `validate:"required"`
	Password  string    `validate:"required"`
}

type RegisterAdminDTO struct {
	Note     string
	Login    string `validate:"required"`
	Password string `validate:"required"`
	Token    string `validate:"required"`
}

type RegisterWaiterDTO struct {
	Login     string `validate:"required"`
	Password  string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Token     string `validate:"required"`
}

type CustomerRegisteredDTO struct {
	CustomerID string `json:"customer_id"`
}
