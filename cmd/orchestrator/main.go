package main

import (
	"finalProject/internal/orchestrator/server"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app := server.New()
	app.RunServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	// Graceful shutdown
	app.Stop()
}
