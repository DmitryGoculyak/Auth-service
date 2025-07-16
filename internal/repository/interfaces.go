package repository

import (
	"Auth-service/internal/entity"
	"context"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(tx *sqlx.Tx, user *entity.User) (*entity.User, error)
	SaveEmail(tx *sqlx.Tx, user *entity.UserEmail) (*entity.UserEmail, error)
	SavePassword(tx *sqlx.Tx, user *entity.UserPassword) (*entity.UserPassword, error)
	FindUserEmail(ctx context.Context, email string) (string, string, error)
	DeleteAllUsers(ctx context.Context) error
}
