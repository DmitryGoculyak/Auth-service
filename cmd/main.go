package main

import (
	"Auth-service/internal/container"
	"Auth-service/pkg/jwt"
	"log"
)

func main() {
	if err := jwt.InitSecret(); err != nil {
		log.Fatal("", err)
		return
	}

	app := container.Build()

	app.Run()
}
