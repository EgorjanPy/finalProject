package agent

import (
	"context"
	pb "finalProject/proto"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

type Application struct {
	Port           string
	ComputingPower int
	stopChan       chan struct{} // Канал для остановки горутин
	wg             sync.WaitGroup
}

func New(port string, computing_power int) *Application {
	return &Application{
		Port:           port,
		ComputingPower: computing_power,
		stopChan:       make(chan struct{}),
	}
}

func (a *Application) StartAgent() {
	defer a.wg.Done() // Уменьшаем счетчик WaitGroup при завершении

	host := "localhost"
	port := "5000"
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		return
	}
	//defer conn.Close()

	grpcClient := pb.NewCalcServiceClient(conn)
	ctx := context.Background()

	for {
		select {
		case <-a.stopChan: // Проверяем сигнал остановки
			log.Println("Stopping agent...")
			return
		default:
			task, err := grpcClient.GetTask(ctx, &pb.GetTaskRequest{})
			if err != nil {
				time.Sleep(time.Second * 2) // Пауза перед повторной попыткой
				continue
			}

			var res float32
			switch task.Operation {
			case "+":
				res = task.Arg1 + task.Arg2
				log.Printf("%d %f + %f = %f", task.Id, task.Arg1, task.Arg2, res)
				time.Sleep(time.Duration(task.OperationTime))
			case "-":
				res = task.Arg1 - task.Arg2
				log.Printf("%d %f - %f = %f", task.Id, task.Arg1, task.Arg2, res)
				time.Sleep(time.Duration(task.OperationTime))

			case "*":
				res = task.Arg1 * task.Arg2
				log.Printf("%d %f * %f = %f", task.Id, task.Arg1, task.Arg2, res)
				time.Sleep(time.Duration(task.OperationTime))

			case "/":
				if task.Arg2 == 0 {
					log.Println("division by zero")
					continue
				}
				res = task.Arg1 / task.Arg2
				log.Printf("%d %f / %f = %f", task.Id, task.Arg1, task.Arg2, res)
				time.Sleep(time.Duration(task.OperationTime))
			}

			_, err = grpcClient.SetTask(ctx, &pb.SetTaskRequest{
				Id:     task.Id,
				Result: res,
			})
			if err != nil {
				log.Println("failed to set task result: ", err)
				continue
			}
			continue
		}
	}
}

func (a *Application) StartApp() {
	// Запускаем указанное количество горутин
	for i := 0; i < a.ComputingPower; i++ {
		a.wg.Add(1)
		go a.StartAgent()
	}
}

func (a *Application) Stop() {
	close(a.stopChan) // Посылаем сигнал остановки всем горутинам
	a.wg.Wait()       // Ждем завершения всех горутин
	log.Println("All agents stopped")
}
