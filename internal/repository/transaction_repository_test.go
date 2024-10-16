package repository

import (
	"context"
	"fmt"
	"my-lime/pkg/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTransactionRepositoryPostgres(sqlxDB)

	rows := sqlmock.NewRows([]string{"id", "transaction_hash", "transaction_status", "block_hash", "block_number", "from_address", "to_address", "contract_address", "logs_count", "input", "value"}).
		AddRow("1", "0x123", 1, "0xabc", 123, "0xfrom", "0xto", "0xcontract", 1, "input", "value")
	query := fmt.Sprintf(`^SELECT (.+) FROM %s`, transactionsTable)
	mock.ExpectQuery(query).WillReturnRows(rows)

	transactions, err := repo.GetAllTransactions(context.Background())

	fmt.Printf("transactions: %v\n", transactions)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, "0x123", transactions[0].TransactionHash)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionsByHashes(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTransactionRepositoryPostgres(sqlxDB)

	txHashes := []string{"0x123", "0x456"}

	rows := sqlmock.NewRows([]string{"id", "transaction_hash", "transaction_status", "block_hash", "block_number", "from_address", "to_address", "contract_address", "logs_count", "input", "value"}).
		AddRow("1", "0x123", 1, "0xabc", 123, "0xfrom", "0xto", "0xcontract", 1, "input1", "value1").
		AddRow("2", "0x456", 1, "0xdef", 456, "0xfrom2", "0xto2", "0xcontract2", 2, "input2", "value2")

	query := fmt.Sprintf(`SELECT \* FROM %s WHERE transaction_hash = ANY\(\$1\)`, transactionsTable)
	mock.ExpectQuery(query).WithArgs(pq.Array(txHashes)).WillReturnRows(rows)

	transactions, err := repo.GetTransactionsByHashes(context.Background(), txHashes)
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, "0x123", transactions[0].TransactionHash)
	assert.Equal(t, "0x456", transactions[1].TransactionHash)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTransactionRepositoryPostgres(sqlxDB)

	tx := models.Transaction{
		TransactionHash:   "0x123",
		TransactionStatus: 0,
		BlockHash:         "",
		BlockNumber:       0,
		From:              "",
		To:                nil,
		ContractAddress:   nil,
		LogsCount:         0,
		Input:             "",
		Value:             "",
	}

	mock.ExpectBegin()
	query := fmt.Sprintf(`INSERT INTO %s \(transaction_hash, transaction_status, block_hash, block_number, from_address, to_address, contract_address, logs_count, input, value\) VALUES \(\?, \?, \?, \?, \?, \?, \?, \?, \?, \?\)`, transactionsTable)
	mock.ExpectExec(query).
		WithArgs(
			tx.TransactionHash,
			tx.TransactionStatus,
			tx.BlockHash,
			tx.BlockNumber,
			tx.From,
			tx.To,
			tx.ContractAddress,
			tx.LogsCount,
			tx.Input,
			tx.Value,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.SaveTransactions(context.Background(), []models.Transaction{tx})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
