package main

import (
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/server"
)

func main() {
	cfg := config.MustLoad()
	app := server.New(cfg.Port)
	// app.Run()545
	app.RunServer()
}
