package main

import (
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app := server.New(config.Cfg.Port)
	// app.Run()
	app.RunServer()

}
