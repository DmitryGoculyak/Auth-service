package pgsql

import (
	"Auth-service/internal/entity"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	db *sqlx.DB
}

func AuthRepoConstructor(
	db *sqlx.DB,
) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) CreateUser(tx *sqlx.Tx, user *entity.User) (*entity.User, error) {

	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	var createUser entity.User
	err := tx.Get(&createUser, "INSERT INTO users (id, full_name) VALUES($1,$2) RETURNING id, full_name, created_at",
		user.ID, user.FullName)
	if err != nil {
		return nil, fmt.Errorf("create user error: %v", err)
	}

	return &createUser, nil
}

func (r *AuthRepo) SaveEmail(tx *sqlx.Tx, user *entity.UserEmail) (*entity.UserEmail, error) {
	var saveEmail entity.UserEmail
	err := tx.Get(&saveEmail, "INSERT INTO emails(user_id, email) VALUES($1, $2) RETURNING id, user_id, email",
		user.UserID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("save email error: %v", err)
	}

	return &saveEmail, nil
}

func (r *AuthRepo) SavePassword(tx *sqlx.Tx, user *entity.UserPassword) (*entity.UserPassword, error) {
	var savePassword entity.UserPassword
	err := tx.Get(&savePassword, "INSERT INTO passwords(user_id, hash) VALUES($1, $2) RETURNING id, user_id, hash, created_at",
		user.UserID, user.Password)
	if err != nil {
		return nil, fmt.Errorf("save password error: %v", err)
	}

	return &savePassword, nil
}

func (r *AuthRepo) FindUserEmail(ctx context.Context, email string) (string, string, error) {
	var userId, hash string
	err := r.db.QueryRowContext(ctx, `
        SELECT u.id, p.hash FROM users u
		JOIN emails e ON u.id = e.user_id
		JOIN passwords p ON u.id = p.user_id
		WHERE e.email = $1
		`, email).Scan(&userId, &hash)
	return userId, hash, err
}

func (r *AuthRepo) DeleteAllUsers(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users")
	return err
}
