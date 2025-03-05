package main

import (
	"finalProject/internal/agent"
	"finalProject/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app := agent.New(cfg.Port, cfg.ComputingPower)
	app.StartApp()
}
