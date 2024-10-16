package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (api *Api) optionalUserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	jwtToken, err := api.extractBearerToken(header)
	if err != nil {
		c.Next()
		return
	}

	userId, err := api.service.Authorization.VerifyToken(c.Request.Context(), jwtToken)
	if err != nil {
		c.Next()
		return
	}

	c.Set(userCtx, userId)
}

func (api *Api) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	jwtToken, err := api.extractBearerToken(header)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	userId, err := api.service.Authorization.VerifyToken(context.Background(), jwtToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	c.Set(userCtx, userId)
}

func (api *Api) extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}
	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
