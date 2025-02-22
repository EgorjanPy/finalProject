package server

import (
	"finalProject/internal/orchestrator/handlers"
	"finalProject/internal/orchestrator/middleware"
	"net/http"
)

type Application struct {
	port string
}

func New(port string) *Application {
	return &Application{port: port}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", middleware.CalcLogger(middleware.CalculateValidation(handlers.CalculateHandler)))
	// http.HandleFunc("/api/v1/expressions", )
	// http.HandleFunc("/api/v1/expressions/:id", )
	// http.HandleFunc("/api/internal/task", ) GET Для получения тасков
	// http.HandleFunc("/api/internal/task", ) POST Для приёма результата таски
	return http.ListenAndServe(a.port, nil)
}
