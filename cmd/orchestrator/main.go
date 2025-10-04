package main

import (
	"finalProject/internal/orchestrator/server"
	"os"
	"os/signal"
)

func main() {
	app := server.New()
	app.StartServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	// Graceful shutdown
	app.Stop()
}
