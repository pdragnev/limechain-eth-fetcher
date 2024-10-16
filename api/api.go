package api

import (
	"fmt"
	"log"
	"my-lime/internal/service"

	"github.com/gin-gonic/gin"
)

type Api struct {
	service *service.Service
}

func NewApi(service *service.Service) *Api {
	return &Api{service: service}
}

func (api *Api) NewRouter(services *service.Service) *gin.Engine {
	r := gin.Default()

	r.SetTrustedProxies(nil)

	transactionAPI := NewTransactionApi(services)
	authAPI := NewAuthApi(services)

	protectedOptional := r.Group("/")
	protectedOptional.Use(api.optionalUserIdentity)
	protectedOptional.GET("/lime/eth", transactionAPI.GetTransactionsByHash)

	protected := r.Group("/")
	protected.Use(api.userIdentity)
	protected.GET("/lime/my", transactionAPI.GetCachedUserTransactions)

	r.GET("/lime/all", transactionAPI.GetAllTransactions)
	r.POST("/lime/authenticate", authAPI.Authenticate)

	return r
}

func (api *Api) StartServer(r *gin.Engine, apiPort int) {
	log.Fatal(r.Run(fmt.Sprintf(":%d", apiPort)))
}
