package domain

import "time"

type RegisterCustomerDTO struct {
	Email     string    `validate:"required,email"`
	Name      string    `validate:"required"`
	Birthdate time.Time `validate:"required"`
	Password  []byte    `validate:"required"`
}

type RegisterAdminDTO struct {
	Note     string
	Login    string `validate:"required"`
	Password []byte `validate:"required"`
	Token    string `validate:"required"`
}

type RegisterWaiterDTO struct {
	Login     string `validate:"required"`
	Password  []byte `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Token     string `validate:"required"`
}

type LoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type TokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
