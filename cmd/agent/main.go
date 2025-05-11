package main

import (
	"context"
	"finalProject/internal/agent"
	"finalProject/internal/config"
	pb "finalProject/proto"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func main() {
	cfg := config.MustLoad()
	app := agent.New(cfg.Port, cfg.ComputingPower)
	app.StartApp()
}
func main() {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()
	/// ..будет продолжение
	grpcClient := pb.NewCalcServiceClient(conn)
	add, err := grpcClient.Add(context.TODO(), &pb.ExprRequest{
		Val1: 10.1,
		Val2: 20.5,
	})

	if err != nil {
		log.Println("failed invoking Area: ", err)
	}

	sub, err := grpcClient.Sub(context.TODO(), &pb.ExprRequest{
		Val1: 10.1,
		Val2: 20.5,
	})
	mul, err := grpcClient.Sub(context.TODO(), &pb.ExprRequest{
		Val1: 10.1,
		Val2: 20.5,
	})
	diff, err := grpcClient.Sub(context.TODO(), &pb.ExprRequest{
		Val1: 10.1,
		Val2: 20.5,
	})
	if err != nil {
		log.Println("failed invoking Area: ", err)
	}

	fmt.Println("Add: ", add.Result)
	fmt.Println("Sub: ", sub.Result)
	fmt.Println("Mul: ", mul.Result)
	fmt.Println("Diff", diff.Result)
}
