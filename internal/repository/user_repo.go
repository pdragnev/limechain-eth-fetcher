package repository

import (
	"context"
	"fmt"
	"my-lime/pkg/models"

	"github.com/jmoiron/sqlx"
)

type UserRepositoryPostgres struct {
	db *sqlx.DB
}

func NewUserRepositoryPostgres(db *sqlx.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) GetUser(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1", usersTable)
	err := r.db.GetContext(ctx, &user, query, username)

	return user, err
}
