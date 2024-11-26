package usecase

import (
	"log/slog"
)

type CustomerAuthRepo interface{}

type AdminAuthRepo interface{}

type WaiterAuthRepo interface{}

type TokensRepo interface{}

type authUsecase struct {
	customers CustomerAuthRepo
	admins    AdminAuthRepo
	waiters   WaiterAuthRepo
	tokens    TokensRepo
	log       *slog.Logger
}

func NewAuthUsecase(log *slog.Logger, customers CustomerAuthRepo, waiters WaiterAuthRepo, admins AdminAuthRepo, tokens TokensRepo) *authUsecase {
	return &authUsecase{
		customers: customers,
		tokens:    tokens,
		admins:    admins,
		waiters:   waiters,
		log:       log,
	}
}
