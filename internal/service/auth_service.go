package service

import (
	"context"
	"fmt"
	"my-lime/internal/config"
	"my-lime/internal/repository"
	"my-lime/internal/utils"
	"my-lime/pkg/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const (
	tokenTTL = 5 * time.Minute
)

type AuthService struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewAuthService(repo repository.UserRepository, config *config.Config) *AuthService {
	return &AuthService{repo: repo, config: config}
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user %w", err)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid username or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt %w", err)
	}

	return token, nil
}

func (s *AuthService) generateToken(user models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.config.JWT_SECRET))
}

func (s *AuthService) VerifyToken(ctx context.Context, jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT_SECRET), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("failed to get subject from claims: %w", err)
	}
	return subject, nil
}
