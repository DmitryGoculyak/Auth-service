package main

import (
	"Auth-service/internal/container"
)

func main() {
	app := container.Build()

	app.Run()
}
