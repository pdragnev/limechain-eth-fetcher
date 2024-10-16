package service

import (
	"context"
	"fmt"
	"my-lime/pkg/models"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionsByHashes(ctx context.Context, txHashes []string) ([]models.Transaction, error) {
	args := m.Called(ctx, txHashes)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) SaveTransactions(ctx context.Context, txs []models.Transaction) error {
	args := m.Called(ctx, txs)
	return args.Error(0)
}

func (m *MockTransactionRepository) SaveTransaction(ctx context.Context, tx models.Transaction) (string, error) {
	args := m.Called(ctx, tx)
	return args.String(0), args.Error(1)
}

func TestGetAllTransactions(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockRepo.On("GetAllTransactions", mock.Anything).Return([]models.Transaction{}, nil)

	service := NewTransactionService(mockRepo, nil, nil)

	transactions, err := service.GetAllTransactions(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, len(transactions), 0)
	mockRepo.AssertExpectations(t)
}

func TestGetTransactionsByHashes(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	txHashes := []string{"0x123", "0x456"}
	mockRepo.On("GetTransactionsByHashes", mock.Anything, txHashes).Return([]models.Transaction{
		{TransactionHash: "0x123"},
		{TransactionHash: "0x456"},
	}, nil)

	service := NewTransactionService(mockRepo, nil, nil)

	transactions, err := service.GetTransactionsByHashes(context.Background(), txHashes)
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, len(transactions), 2)
	mockRepo.AssertExpectations(t)
}

func TestCacheUserTransactions(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	userId := "test-user"
	key := fmt.Sprintf("userid:%s:txhashes", userId)
	txHashes := []string{"0x123", "0x456"}
	mockRedis.ExpectSAdd(key, txHashes).SetVal(int64(2))
	service := NewTransactionService(nil, redisClient, nil)

	err := service.CacheUserTransactions(context.Background(), userId, []models.Transaction{
		{TransactionHash: "0x123"},
		{TransactionHash: "0x456"},
	})
	assert.NoError(t, err)
	err = mockRedis.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetCachedUserTransactions(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	userId := "test-user"
	key := fmt.Sprintf("userid:%s:txhashes", userId)
	txHashes := []string{"0x123", "0x456"}
	mockRedis.ExpectSMembers(key).SetVal(txHashes)

	mockRepo := new(MockTransactionRepository)
	mockRepo.On("GetTransactionsByHashes", mock.Anything, txHashes).Return([]models.Transaction{
		{TransactionHash: "0x123"},
		{TransactionHash: "0x456"},
	}, nil)

	service := NewTransactionService(mockRepo, redisClient, nil)

	txs, err := service.GetCachedUserTransactions(context.Background(), userId)
	assert.NoError(t, err)
	assert.Equal(t, len(txs), len(txHashes))
	err = mockRedis.ExpectationsWereMet()
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAndSaveTransactionByHashes(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	txHashes := []string{"0x123", "0x456"}
	mockRepo.On("GetTransactionsByHashes", mock.Anything, txHashes).Return([]models.Transaction{
		{TransactionHash: "0x123"},
		{TransactionHash: "0x456"},
	}, nil)
	service := NewTransactionService(mockRepo, nil, nil)

	txs, err := service.GetAndSaveTransactionByHashes(context.Background(), txHashes)
	assert.NoError(t, err)
	assert.Equal(t, len(txs), len(txHashes))
	mockRepo.AssertExpectations(t)
}
