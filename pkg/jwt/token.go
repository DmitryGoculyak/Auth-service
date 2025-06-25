package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type JWTConfig struct {
	SecretKey string
}

var secret string

func InitSecret() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("JWT_SECRET is not set")
	}
	return nil
}

func GenerateJWTToken(userID string) (string, error) {
	if secret == "" {
		return "", errors.New("secret key is empty")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
