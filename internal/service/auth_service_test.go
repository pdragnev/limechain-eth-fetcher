package service

import (
	"context"
	"my-lime/internal/config"
	"my-lime/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSecret = "test_secret"

var testConfig = &config.Config{JWT_SECRET: testSecret}

func TestGenerateToken(t *testing.T) {
	service := NewAuthService(nil, testConfig)

	user := models.User{ID: "test-user"}
	token, err := service.generateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {
	service := NewAuthService(nil, testConfig)

	user := models.User{ID: "test-user"}
	tokenString, err := service.generateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	userId, err := service.VerifyToken(context.Background(), tokenString)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userId)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	service := NewAuthService(nil, testConfig)

	invalidToken := "invalid.token.string"
	userId, err := service.VerifyToken(context.Background(), invalidToken)
	assert.Error(t, err)
	assert.Empty(t, userId)
}
