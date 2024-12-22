package main

import (
	"finalProject/internal/application"
)

func main() {
	app := application.New()
	// app.Run()
	app.RunServer()
}
