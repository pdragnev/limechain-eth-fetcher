package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	transactionsTable = "transactions"
	usersTable        = "users"
)

func NewPostgresDB(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
