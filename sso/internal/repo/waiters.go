package repo

import "github.com/jmoiron/sqlx"

type waiterRepo struct {
	db *sqlx.DB
}

func NewWaiterRepo(db *sqlx.DB) *waiterRepo {
	return &waiterRepo{db: db}
}
