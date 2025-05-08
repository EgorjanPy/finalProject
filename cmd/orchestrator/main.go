package main

import (
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/server"
	_ "github.com/mattn/go-sqlite3"
)

//	func CreateTables(db *sql.DB) error {
//		const (
//			usersTable = `
//		CREATE TABLE IF NOT EXISTS users(
//			id INTEGER PRIMARY KEY AUTOINCREMENT,
//			login TEXT NOT NULL UNIQUE,
//			password TEXT NOT NULL
//		);`
//
//			expressionsTable = `
//		CREATE TABLE IF NOT EXISTS expressions(
//			id INTEGER PRIMARY KEY AUTOINCREMENT,
//			expression TEXT NOT NULL,
//			user_id INTEGER NOT NULL,
//			answer TEXT,
//			status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'calculating', 'completed', 'failed')),
//			FOREIGN KEY (user_id)  REFERENCES user(id) ON DELETE CASCADE
//		);`
//		)
//		if _, err := db.Exec(usersTable); err != nil {
//			return err
//		}
//		if _, err := db.Exec(expressionsTable); err != nil {
//			return err
//		}
//		return nil
//	}
func main() {
	cfg := config.MustLoad()
	app := server.New(cfg.Port)
	// app.Run()
	app.RunServer()

}
