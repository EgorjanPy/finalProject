package server

import (
	"finalProject/internal/orchestrator/handlers"
	"finalProject/internal/orchestrator/middleware"
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

	r.HandleFunc("/api/v1/calculate", middleware.LoggerMiddleware(handlers.CalculateHandler))
	r.HandleFunc("/api/v1/expressions", middleware.LoggerMiddleware(handlers.ExpressionsHandler))
	r.HandleFunc("/api/v1/expressions/{id}", middleware.LoggerMiddleware(handlers.GetExpressionByIdHandler))
	r.HandleFunc("/internal/task", handlers.GetSetTask)
	http.Handle("/", r)
	fmt.Printf("Сервер удачно запущен на http://localhost%s", a.port)
	fmt.Println()
	return http.ListenAndServe(a.port, nil)
}
