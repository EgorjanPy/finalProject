package main

import (
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()
	app := server.New(cfg.Port)
	// app.Run()
	app.RunServer()

}
