package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(tx *sqlx.Tx, fullName string) (string, error)
	SaveEmail(tx *sqlx.Tx, userID, email string) error
	SavePassword(tx *sqlx.Tx, userID, hash string) error
	FindUserEmail(ctx context.Context, email string) (string, string, error)
	DeleteAllUsers(ctx context.Context) error
}
