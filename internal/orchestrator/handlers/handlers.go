package handlers

import (
	"encoding/json"
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/logic"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/mux"
)

var cfg = config.MustLoad()

type ExpressionsResponse struct {
	Expressions []logic.Expression
}

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
}
func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	response := ExpressionsResponse{Expressions: logic.Expressions.GetExpressions()}
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

func isValidExpression(expression string) bool {
	cleanedExpression := removeSpaces(expression)
	validPattern := `^[0-9+\-*/()]+$`
	matched, err := regexp.MatchString(validPattern, cleanedExpression)
	if err != nil || !matched {
		return false
	}
	if !areParenthesesBalanced(cleanedExpression) {
		return false
	}
	if strings.Contains("+=-/*:", string(expression[0])) || strings.Contains("+=-/*:", string(expression[len(expression)-1])) {
		return false
	}
	for i := 1; i < len(expression)-1; i++ {
		if rune(expression[i-1]) == '/' && rune(expression[i-1]) == '0' {
			return false
		}
		if isSign(rune(expression[i-1])) && isSign(rune(expression[i])) {
			return false
		}
	}
	return true
}
func removeSpaces(s string) string {
	var result []rune
	for _, r := range s {
		if !unicode.IsSpace(r) {
			result = append(result, r)
		}
	}
	return string(result)
}

func areParenthesesBalanced(expression string) bool {
	var stack []rune
	for _, char := range expression {
		if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			if len(stack) == 0 {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("cant read body :(")
		w.WriteHeader(500)
	}
	defer r.Body.Close()
	var request CalculateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal("cant unmarsahl body :(")
		w.WriteHeader(500)
	}
	if !isValidExpression(request.Expression) {
		w.WriteHeader(422)
		return
	}
	id := logic.NewEx(request.Expression)
	response := CalculateResponse{Id: id}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("cant marsahl response :(")
		w.WriteHeader(500)
	}
	w.WriteHeader(201)
	w.Write(jsonBytes)
}

type GetExpressionResponse struct {
	Expression logic.Expression
}

func GetExpressionByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	ex, err := logic.Expressions.GetExpressionById(id)
	if err != nil {
		w.WriteHeader(404)
	}
	response := GetExpressionResponse{ex}
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

func isValidResult(result string) bool {
	cleanedExpression := removeSpaces(result)
	validPattern := `[0-9]`
	matched, err := regexp.MatchString(validPattern, cleanedExpression)
	if err != nil || !matched {
		return false
	}
	return true

}
func GetSetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("cant read body :(")
			w.WriteHeader(500)
			return
		}
		defer r.Body.Close()
		var request GetSetTaskRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Fatal("cant unmarsahl body :(")
			w.WriteHeader(500)
			return
		}

		if logic.Tasks.GetLen() <= request.Id {
			w.WriteHeader(404)
			return
		}
		logic.Results.SetResult(request.Id, request.Result)
		return
	}
	if r.Method == http.MethodGet {
		id := logic.Results.GetLen()
		task, err := logic.Tasks.GetTaskById(id)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		response := GetSetTaskResponse{Id: id, Arg1: task.Arg1, Arg2: task.Arg2, Operation: task.Operation}
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

type RegisterRequest struct {
	Login    string
	Password string
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("cant read body :(")
			w.WriteHeader(500)
			return
		}
		defer r.Body.Close()
		var request RegisterRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Fatal("cant unmarsahl body :(")
			w.WriteHeader(500)
			return
		}
		login := request.Login
		password := request.Password
		// Проверка на занятность логина
		// Добавление в БД
	}
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("cant read body :(")
			w.WriteHeader(500)
			return
		}
		defer r.Body.Close()
		var request RegisterRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Fatal("cant unmarsahl body :(")
			w.WriteHeader(500)
			return
		}
		login := request.Login
		password := request.Password
		// Проверка пароля по логину в БД
		// Если всё гуд то генерируем токен
		t, rt := logic.GengerateJWT(login)
		r.AddCookie(&http.Cookie{
			Name:  "token",
			Value: t,
		})
		r.AddCookie(
			&http.Cookie{
				Name:  "refresh_token",
				Value: rt,
			})
}
