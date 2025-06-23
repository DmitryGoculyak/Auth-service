package db

import (
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	conn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	DB, err = sqlx.Connect("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected!")
}
