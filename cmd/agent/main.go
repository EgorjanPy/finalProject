package main

import (
	"finalProject/internal/agent"
	"finalProject/internal/config"
	"os"
	"os/signal"
)

func main() {
	app := agent.New(config.Cfg.Port, config.Cfg.ComputingPower)
	app.StartApp()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Graceful shutdown
	app.Stop()
}
