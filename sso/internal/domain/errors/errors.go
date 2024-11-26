package errs

import "errors"

var (
	ErrInvalidSecretToken    = errors.New("invalid secret token")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrWaiterAlreadyExists   = errors.New("customer already exists")
	ErrWaiterNotFound        = errors.New("waiter not found")
	ErrAdminAlreadyExists    = errors.New("admin already exists")
	ErrAdminNotFound         = errors.New("admin not found")
)
