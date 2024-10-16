package api

import (
	"my-lime/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthApi struct {
	service *service.Service
}

func NewAuthApi(service *service.Service) *AuthApi {
	return &AuthApi{service: service}
}

func (api *AuthApi) Authenticate(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := api.service.Authorization.Authenticate(c.Request.Context(), credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
