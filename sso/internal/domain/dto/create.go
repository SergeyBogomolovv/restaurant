package dto

import "time"

type CreateCustomerDTO struct {
	Email     string
	Name      string
	Birthdate time.Time
	Password  []byte
}

type CreateAdminDTO struct {
	Login    string
	Password []byte
	Note     string
}

type CreateWaiterDTO struct {
	Login     string
	Password  []byte
	FirstName string
	LastName  string
}
