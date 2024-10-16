package service

import (
	"context"
	"my-lime/internal/config"
	"my-lime/internal/repository"
	"my-lime/pkg/models"

	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}

type Transaction interface {
	GetAndSaveTransactionByHashes(ctx context.Context, txHashes []string) ([]models.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]models.Transaction, error)
	GetTransactionsByHashes(ctx context.Context, txHashes []string) ([]models.Transaction, error)
	CacheUserTransactions(ctx context.Context, userId string, txs []models.Transaction) error
	GetCachedUserTransactions(ctx context.Context, userId string) ([]models.Transaction, error)
}

type Service struct {
	Transaction
	Authorization
}

func NewService(repositories *repository.Repository, redisClient *redis.Client, config *config.Config, ethClient EthereumClient) *Service {
	return &Service{
		Transaction:   NewTransactionService(repositories.TransactionRepository, redisClient, ethClient),
		Authorization: NewAuthService(repositories.UserRepository, config),
	}
}
