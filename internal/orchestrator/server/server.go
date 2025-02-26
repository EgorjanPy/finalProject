package server

import (
	"finalProject/internal/orchestrator/handlers"
	"fmt"
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

	r.HandleFunc("/internal/task", handlers.GetTask) // GET Для получения тасков
	// http.HandleFunc("/api/internal/task", ) POST Для приёма результата таски
	http.Handle("/", r)
	fmt.Printf("Сервер удачно запущен на http://localhost%s", a.port)
	return http.ListenAndServe(a.port, nil)
}
