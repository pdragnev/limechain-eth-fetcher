package utils

import (
	"encoding/hex"
	"my-lime/pkg/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/bcrypt"
)

func ToTransaction(receipt *types.Receipt, tx *types.Transaction, from common.Address) (*models.Transaction, error) {
	return &models.Transaction{
		BlockHash:         receipt.BlockHash.Hex(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		TransactionHash:   tx.Hash().Hex(),
		From:              from.Hex(),
		To:                adrToHexPointer(tx.To()),
		ContractAddress:   adrToHexPointer(&receipt.ContractAddress),
		LogsCount:         len(receipt.Logs),
		TransactionStatus: receipt.Status,
		Input:             hex.EncodeToString(tx.Data()),
		Value:             tx.Value().String(),
	}, nil
}

func adrToHexPointer(adr *common.Address) *string {
	if adr == nil || adr.Hex() == "0x0000000000000000000000000000000000000000" {
		return nil
	}
	adrHex := adr.Hex()
	return &adrHex
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
