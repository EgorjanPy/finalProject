package server

import (
	"context"
	"finalProject/internal/orchestrator/handlers"
	"finalProject/internal/orchestrator/logic"
	"finalProject/internal/orchestrator/middleware"
	pb "finalProject/proto"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

type Application struct {
	port string
}
type Server struct {
	pb.CalcServiceServer // Сервис из сгенерированного пакета
}

func NewServer() *Server {
	return &Server{}
}
func New(port string) *Application {
	return &Application{port: port}
}

func (s *Server) GetTask(ctx context.Context, in *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	//fmt.Println("GET")
	id := logic.Results.GetLen()
	task, err := logic.Tasks.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	return &pb.GetTaskResponse{Id: int32(id), Arg1: float32(task.Arg1), Arg2: float32(task.Arg2), Operation: task.Operation}, nil
}

func (s *Server) SetTask(ctx context.Context, in *pb.SetTaskRequest) (*pb.SetTaskResponse, error) {
	//fmt.Println("SET")
	logic.Results.SetResult(int(in.Id), float64(in.Result))
	return &pb.SetTaskResponse{}, nil
}
func (a *Application) RunServer() (error, error) {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	calcServiceServer := NewServer()
	pb.RegisterCalcServiceServer(grpcServer, calcServiceServer)
	go func() {
		log.Println("tcp listener started at port: ", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", middleware.LoggerMiddleware(middleware.ProtectedHandler(handlers.CalculateHandler)))
	r.HandleFunc("/api/v1/expressions", middleware.LoggerMiddleware(middleware.ProtectedHandler(handlers.ExpressionsHandler)))
	r.HandleFunc("/api/v1/expressions/{id}", middleware.LoggerMiddleware(middleware.ProtectedHandler(handlers.GetExpressionByIdHandler)))
	r.HandleFunc("/api/v1/register", handlers.RegisterHandler)
	r.HandleFunc("/api/v1/login", handlers.LoginHandler)
	http.Handle("/", r)
	log.Println()
	go func() {
		log.Printf("Сервер удачно запущен на http://localhost%s", a.port)
		if err := http.ListenAndServe(a.port, nil); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Бесконечный цикл, чтобы main не завершился
	select {}
}
