package repository

import (
	"context"
	"my-lime/pkg/models"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	SaveTransactions(ctx context.Context, txs []models.Transaction) error
	GetTransactionsByHashes(ctx context.Context, hashes []string) ([]models.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]models.Transaction, error)
}

type UserRepository interface {
	GetUser(ctx context.Context, username string) (models.User, error)
}

type Repository struct {
	TransactionRepository
	UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TransactionRepository: NewTransactionRepositoryPostgres(db),
		UserRepository:        NewUserRepositoryPostgres(db),
	}
}
