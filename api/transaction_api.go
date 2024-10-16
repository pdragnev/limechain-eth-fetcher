package api

import (
	"my-lime/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionApi struct {
	service *service.Service
}

func NewTransactionApi(service *service.Service) *TransactionApi {
	return &TransactionApi{service: service}
}

func (api *TransactionApi) GetTransactionsByHash(c *gin.Context) {
	userId, exist := c.Get("userId")
	txHashStrs := c.QueryArray("transactionHashes")

	if len(txHashStrs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transactionHashes query parameter is required"})
		return
	}

	transactions, err := api.service.Transaction.GetAndSaveTransactionByHashes(c.Request.Context(), txHashStrs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exist {
		err := api.service.Transaction.CacheUserTransactions(c.Request.Context(), userId.(string), transactions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, transactions)
}

func (api *TransactionApi) GetAllTransactions(c *gin.Context) {
	transactions, err := api.service.Transaction.GetAllTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (api *TransactionApi) GetCachedUserTransactions(c *gin.Context) {
	userId, exist := c.Get("userId")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token id not parsed from middleware"})
		return
	}
	transactions, err := api.service.GetCachedUserTransactions(c.Request.Context(), userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)

}
