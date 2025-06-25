package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(cfg *DBConfig) (*sqlx.DB, error) {
	var err error
	var db *sqlx.DB

	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)
	db, err = sqlx.Connect("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("[DATABASE] Successfully connected!")
	return db, nil
}
