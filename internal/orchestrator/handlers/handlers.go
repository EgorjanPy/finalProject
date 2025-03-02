package handlers

import (
	"encoding/json"
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/logic"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var cfg = config.MustLoad()

type ExpressionsResponse struct {
	Expressions []logic.Expression
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	response := ExpressionsResponse{Expressions: logic.Expressions}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("cant marsahl response :(")
		w.WriteHeader(500)
		return
	}
	w.Write(jsonBytes)
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}
type CalculateResponse struct {
	Id int `json:"id"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("cant read body :(")
		w.WriteHeader(422)
		return
	}
	// fmt.Println("Ok 1")
	defer r.Body.Close()
	var request CalculateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal("cant unmarsahl body :(")
		w.WriteHeader(422)
		return
	}
	// fmt.Println("Ok 2")
	id := logic.NewEx(request.Expression)
	response := CalculateResponse{Id: id}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("cant marsahl response :(")
		w.WriteHeader(500)
		return
	}
	// fmt.Println("Ok 3")

	w.Write(jsonBytes)
	// log.Printf("%d", id)
}

type GetExpressionResponse struct {
	Expression logic.Expression
}

func GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	response := GetExpressionResponse{logic.Expressions[id]}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("cant marsahl response :(")
		w.WriteHeader(500)
		return
	}
	w.Write(jsonBytes)
}

type GetSetTaskResponse struct {
	Id             int           `json:"id"`
	Arg1           float64       `json:"arg1"`
	Arg2           float64       `json:"arg2"`
	Operation      string        `json:"operation"`
	Operation_time time.Duration `json:"operation_time"`
}
type GetSetTaskRequest struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

func GetSetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("cant read body :(")
			w.WriteHeader(422)
			return
		}
		defer r.Body.Close()
		var request GetSetTaskRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Fatal("cant unmarsahl body :(")
			w.WriteHeader(422)
			return
		}
		logic.Results.SetResult(request.Id, request.Result)
		fmt.Println(logic.Results.Results)
		fmt.Println(logic.Tasks.Tasks)
		return
	}
	if r.Method == http.MethodGet {
		i := len(logic.Results.Results)

		task := logic.Tasks.Tasks[i]
		response := GetSetTaskResponse{Id: task.Id, Arg1: task.Arg1, Arg2: task.Arg2, Operation: task.Operation}
		switch task.Operation {
		case "+":
			response.Operation_time = cfg.TimeAddMs
		case "-":
			response.Operation_time = cfg.TimeDivMs
		case "*":
			response.Operation_time = cfg.TimeSubMs
		case "/":
			response.Operation_time = cfg.TimeMulMs
		}
		jsonRes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(jsonRes)
		return
	}

}
