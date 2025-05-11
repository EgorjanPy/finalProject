package server

import (
	"context"
	"finalProject/internal/orchestrator/handlers"
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
func (s *Server) Add(ctx context.Context, in *pb.ExprRequest) (*pb.AddResponse, error) {
	log.Println("invoked add: ", in)
	// вычислим площадь и вернём ответ
	return &pb.AddResponse{
		Result: in.Val1 + in.Val2,
	}, nil
}
func (s *Server) Sub(ctx context.Context, in *pb.ExprRequest) (*pb.SubResponse, error) {
	log.Println("invoked sub: ", in)
	// вычислим площадь и вернём ответ
	return &pb.SubResponse{
		Result: in.Val1 - in.Val2,
	}, nil
}
func (s *Server) Mul(ctx context.Context, in *pb.ExprRequest) (*pb.MulResponse, error) {
	log.Println("invoked mul: ", in)
	// вычислим площадь и вернём ответ
	return &pb.MulResponse{
		Result: in.Val1 * in.Val2,
	}, nil
}
func (s *Server) Diff(ctx context.Context, in *pb.ExprRequest) (*pb.DiffResponse, error) {
	log.Println("invoked diff: ", in)
	// вычислим площадь и вернём ответ
	return &pb.DiffResponse{
		Result: in.Val1 / in.Val2,
	}, nil
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
	log.Println("tcp listener started at port: ", port)
	grpcServer := grpc.NewServer()
	calcServiceServer := NewServer()
	pb.RegisterCalcServiceServer(grpcServer, calcServiceServer.CalcServiceServer)
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", middleware.LoggerMiddleware(handlers.CalculateHandler))
	r.HandleFunc("/api/v1/expressions", middleware.LoggerMiddleware(handlers.ExpressionsHandler))
	r.HandleFunc("/api/v1/expressions/{id}", middleware.LoggerMiddleware(handlers.GetExpressionByIdHandler))
	r.HandleFunc("/internal/task", handlers.GetSetTask)
	r.HandleFunc("/api/v1/register", handlers.RegisterHandler)
	r.HandleFunc("/api/v1/login", handlers.LoginHandler)
	http.Handle("/", r)
	log.Printf("Сервер удачно запущен на http://localhost%s", a.port)
	log.Println()
	//if err := grpcServer.Serve(lis); err != nil {
	//	log.Println("error serving grpc: ", err)
	//	os.Exit(1)
	//}
	return http.ListenAndServe(a.port, nil), grpcServer.Serve(lis)
}
