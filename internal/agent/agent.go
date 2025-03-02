package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Id     int
	Result float64
}
type Response struct {
	Id             int           `json:"id"`
	Arg1           float64       `json:"arg1"`
	Arg2           float64       `json:"arg2"`
	Operation      string        `json:"operation"`
	Operation_time time.Duration `json:"operation_time"`
}

func main() {
	url := "http://localhost:8080/internal/task"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var response Response
	json.Unmarshal(body, &response)
	fmt.Println(response)
	var res float64
	switch response.Operation {
	case "+":
		res = response.Arg1 + response.Arg2
		time.Sleep(response.Operation_time)
	case "-":
		res = response.Arg1 - response.Arg2
		time.Sleep(response.Operation_time)
	case "*":
		res = response.Arg1 * response.Arg2
		time.Sleep(response.Operation_time)
	case "/":
		res = response.Arg1 / response.Arg2
		time.Sleep(response.Operation_time)
	}
	request := Request{}
	request.Id = response.Id
	request.Result = res
	jsonData, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response2, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer response2.Body.Close()
	// log.Println(response2)
}
