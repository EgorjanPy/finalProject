package server

import (
	"context"
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/handlers"
	"finalProject/internal/orchestrator/logic"
	"finalProject/internal/orchestrator/middleware"
	"finalProject/internal/storage"
	pb "finalProject/proto"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type Application struct {
	port     string
	tcpPort  string
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func New() *Application {
	return &Application{
		port:     config.Cfg.Port,
		tcpPort:  config.Cfg.TCPPort,
		stopChan: make(chan struct{})}
}

type Server struct {
	pb.CalcServiceServer // Сервис из сгенерированного пакета
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetTask(ctx context.Context, in *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	id := logic.Results.GetLen()
	task, err := logic.Tasks.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	return &pb.GetTaskResponse{Id: int32(id), Arg1: float32(task.Arg1), Arg2: float32(task.Arg2), Operation: task.Operation}, nil
}
func (s *Server) SetTask(ctx context.Context, in *pb.SetTaskRequest) (*pb.SetTaskResponse, error) {
	if in.Error {
	}
	logic.Results.SetResult(int(in.Id), float64(in.Result))
	return &pb.SetTaskResponse{}, nil
}

func (a *Application) RunServer() {
	host := "localhost"
	tcpPort := a.tcpPort

	addr := fmt.Sprintf("%s%s", host, tcpPort)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу
	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	calcServiceServer := NewServer()
	pb.RegisterCalcServiceServer(grpcServer, calcServiceServer)
	go func() {
		log.Println("слушатель tcp запущен на порту: ", tcpPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	r := mux.NewRouter()

	// FRONTEND
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", middleware.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./static/index.html") }))
	r.HandleFunc("/expressions", middleware.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./static/expressions.html") }))
	r.HandleFunc("/expressions/{id}", middleware.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./static/expressions.html") }))
	r.HandleFunc("/register", middleware.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./static/register.html") }))
	r.HandleFunc("/login", middleware.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./static/login.html") }))

	// API
	r.HandleFunc("/api/v1/calculate", middleware.LoggerMiddleware(middleware.ProtectedHandler(handlers.CalculateHandler)))
	r.HandleFunc("/api/v1/expressions", middleware.LoggerMiddleware(middleware.ProtectedHandler(handlers.ExpressionsHandler)))
	r.HandleFunc("/api/v1/expressions/{expressionID}", middleware.LoggerMiddleware(middleware.ProtectedHandler(handlers.GetExpressionByIdHandler)))
	r.HandleFunc("/api/v1/register", middleware.LoggerMiddleware(handlers.RegisterHandler))
	r.HandleFunc("/api/v1/login", middleware.LoggerMiddleware(handlers.LoginHandler))
	http.Handle("/", r)

	// Проверка есть ли нерешенные выражения в бд, если да, то добавляем их в очередь
	expressions, err := storage.DataBase.GetUncompletedExpressions()
	if err != nil {
		log.Println("can't get uncompleted expressions from database")
	} else {
		go func() {
			for _, ex := range expressions {
				logic.NewExpression(int(ex.ID), ex.Expression, ex.UserID)
			}
		}()

	}

	go func() {
		log.Printf("Сервер удачно запущен на %s%s", host, a.port)
		if err := http.ListenAndServe(a.port, nil); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Бесконечный цикл, чтобы main не завершился
	select {
	case <-a.stopChan: // Проверяем сигнал остановки
		log.Println("Stopping server...")
		return
	}
}
func (a *Application) Stop() {
	close(a.stopChan) // Посылаем сигнал остановки всем горутинам
	a.wg.Wait()       // Ждем завершения всех горутин
	log.Println("All agents stopped")
}
