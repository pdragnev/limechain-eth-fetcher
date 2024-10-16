package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"my-lime/internal/repository"
	"my-lime/internal/utils"
	"my-lime/pkg/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/redis/go-redis/v9"
)

type TransactionService struct {
	repo        repository.TransactionRepository
	redisClient *redis.Client
	ethClient   EthereumClient
}

func NewTransactionService(repo repository.TransactionRepository, redisClient *redis.Client, ethClient EthereumClient) *TransactionService {
	return &TransactionService{repo: repo, redisClient: redisClient, ethClient: ethClient}
}

func (s *TransactionService) GetAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	return s.repo.GetAllTransactions(ctx)
}

func (s *TransactionService) GetTransactionsByHashes(ctx context.Context, txHashes []string) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByHashes(ctx, txHashes)
}

func (s *TransactionService) CacheUserTransactions(ctx context.Context, userId string, txs []models.Transaction) error {
	var txHashes []string
	for _, tx := range txs {
		txHashes = append(txHashes, tx.TransactionHash)
	}
	key := fmt.Sprintf("userid:%s:txhashes", userId)
	_, err := s.redisClient.SAdd(ctx, key, txHashes).Result()
	return err
}

func (s *TransactionService) GetCachedUserTransactions(ctx context.Context, userId string) ([]models.Transaction, error) {
	key := fmt.Sprintf("userid:%s:txhashes", userId)
	txHashes, err := s.redisClient.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get result from redis: %w", err)
	}

	return s.GetTransactionsByHashes(ctx, txHashes)
}

func (s *TransactionService) GetAndSaveTransactionByHashes(ctx context.Context, txHashes []string) ([]models.Transaction, error) {
	existingTransactions, err := s.GetTransactionsByHashes(ctx, txHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing transaction: %w", err)
	}

	var resultTxs []models.Transaction
	var missingHashes []common.Hash

	existingTransactionsMap := make(map[string]models.Transaction)
	for _, tx := range existingTransactions {
		existingTransactionsMap[tx.TransactionHash] = tx
	}

	for _, txHashStr := range txHashes {
		if tx, found := existingTransactionsMap[txHashStr]; found {
			resultTxs = append(resultTxs, tx)
		} else {
			missingHashes = append(missingHashes, common.HexToHash(txHashStr))
		}
	}

	if len(missingHashes) > 0 {
		fetchedTxs, err := s.fetchAndSaveMissingTransactions(ctx, missingHashes)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch missing transaction: %w", err)
		}
		resultTxs = append(resultTxs, fetchedTxs...)
	}

	return resultTxs, nil
}

func (s *TransactionService) fetchAndSaveMissingTransactions(ctx context.Context, txHashes []common.Hash) ([]models.Transaction, error) {
	fetchedTxs, err := s.fetchTransactionsFromNode(ctx, txHashes)
	if err != nil {
		return nil, err
	}
	if len(fetchedTxs) == 0 {
		return nil, errors.New("failed to fetch transactions from node")
	}
	err = s.repo.SaveTransactions(ctx, fetchedTxs)
	if err != nil {
		return nil, fmt.Errorf("failed to save transactions: %w", err)
	}

	return fetchedTxs, nil
}

func (s *TransactionService) fetchTransactionsFromNode(ctx context.Context, txHashes []common.Hash) ([]models.Transaction, error) {
	var fetchedTxs []models.Transaction

	for _, txHash := range txHashes {
		tx, isPending, err := s.ethClient.TransactionByHash(ctx, txHash)
		if err != nil {
			if err.Error() == "not found" {
				log.Printf("failed to get tx by hash. %s %v", txHash, err)
				continue
			}
			return nil, fmt.Errorf("failed to get transaction by hash %s: %w", txHash.Hex(), err)

		}
		if isPending {
			continue
		}

		receipt, err := s.ethClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction receipt for %s: %w", txHash.Hex(), err)
		}

		chainID, err := s.ethClient.NetworkID(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get network ID: %w", err)
		}

		from, err := types.Sender(types.NewLondonSigner(chainID), tx)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction sender: %w", err)
		}

		resultTx, err := utils.ToTransaction(receipt, tx, from)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction: %w", err)
		}

		fetchedTxs = append(fetchedTxs, *resultTx)
	}
	return fetchedTxs, nil
}
