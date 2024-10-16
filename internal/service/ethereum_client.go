package service

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient interface {
	TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	NetworkID(ctx context.Context) (*big.Int, error)
}

type EthereumClientImpl struct {
	client *ethclient.Client
}

func NewEthereumClient(url string) (EthereumClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthereumClientImpl{client: client}, nil
}

func (e *EthereumClientImpl) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	return e.client.TransactionByHash(ctx, txHash)
}

func (e *EthereumClientImpl) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return e.client.TransactionReceipt(ctx, txHash)
}

func (e *EthereumClientImpl) NetworkID(ctx context.Context) (*big.Int, error) {
	return e.client.NetworkID(ctx)
}
