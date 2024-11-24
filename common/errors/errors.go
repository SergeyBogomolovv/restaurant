package errors

import "errors"

var (
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidToken          = errors.New("invalid token")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrWaiterAlreadyExists   = errors.New("customer already exists")
	ErrWaiterNotFound        = errors.New("waiter not found")
	ErrAdminNotFound         = errors.New("admin not found")
	ErrAdminAlreadyExists    = errors.New("admin already exists")
)
