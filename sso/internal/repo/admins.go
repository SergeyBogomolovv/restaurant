package repo

import "github.com/jmoiron/sqlx"

type adminRepo struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *adminRepo {
	return &adminRepo{db: db}
}
