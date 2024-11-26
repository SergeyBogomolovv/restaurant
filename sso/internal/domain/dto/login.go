package dto

type LoginCustomerDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginEmployeeDTO struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
}
