package main

import (
	"finalProject/internal/agent"
	"os"
	"os/signal"
)

func main() {
	app := agent.New()
	app.StartApp()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	// Graceful shutdown
	app.Stop()
}
