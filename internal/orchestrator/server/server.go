package server

import (
	"finalProject/internal/orchestrator/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type Application struct {
	port string
}

func New(port string) *Application {
	return &Application{port: port}
}

func (a *Application) RunServer() error {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)
	r.HandleFunc("/api/v1/expressions", handlers.ExpressionsHandler)
	r.HandleFunc("/api/v1/expressions/{id}", handlers.GetExpressionHandler)
	http.Handle("/", r)
	// http.HandleFunc("/api/internal/task", ) GET Для получения тасков
	// http.HandleFunc("/api/internal/task", ) POST Для приёма результата таски
	return http.ListenAndServe(a.port, nil)
}
