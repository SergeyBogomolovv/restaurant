package usecase

import "log/slog"

type authUsecase struct {
	customers CustomerRepo
	admins    AdminRepo
	waiters   WaiterRepo
	tokens    TokensRepo
	log       *slog.Logger
}

func NewAuthUsecase(log *slog.Logger, customers CustomerRepo, waiters WaiterRepo, admins AdminRepo, tokens TokensRepo) *authUsecase {
	return &authUsecase{
		customers: customers,
		tokens:    tokens,
		admins:    admins,
		waiters:   waiters,
		log:       log,
	}
}
