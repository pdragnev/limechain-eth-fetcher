package repository

import (
	"context"
	"fmt"
	"log"
	"my-lime/pkg/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TransactionRepositoryPostgres struct {
	db *sqlx.DB
}

func NewTransactionRepositoryPostgres(db *sqlx.DB) *TransactionRepositoryPostgres {
	return &TransactionRepositoryPostgres{db: db}
}

func (r *TransactionRepositoryPostgres) SaveTransactions(ctx context.Context, txs []models.Transaction) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (transaction_hash, transaction_status, block_hash, block_number, from_address, to_address, contract_address, logs_count, input, value)
		VALUES (:transaction_hash, :transaction_status, :block_hash, :block_number, :from_address, :to_address, :contract_address, :logs_count, :input, :value)
	`, transactionsTable)

	dbTx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}

	_, err = dbTx.NamedExecContext(ctx, query, txs)
	if err != nil {
		dbTx.Rollback()
		log.Printf("Failed to execute batch insert: %v", err)
		return err
	}

	err = dbTx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (r *TransactionRepositoryPostgres) GetTransactionsByHashes(ctx context.Context, txHashes []string) ([]models.Transaction, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE transaction_hash = ANY($1)`, transactionsTable)
	var transactions []models.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, pq.Array(txHashes))
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepositoryPostgres) GetAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := fmt.Sprintf(`SELECT * FROM %s`, transactionsTable)
	err := r.db.SelectContext(ctx, &transactions, query)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
