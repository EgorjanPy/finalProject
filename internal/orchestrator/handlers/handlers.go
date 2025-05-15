package handlers

import (
	"encoding/json"
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/logic"
	"finalProject/internal/storage"
	"finalProject/internal/storage/sqlite"
	"fmt"
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
	Expressions []sqlite.Expression
}

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
}
func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressions, err := storage.DataBase.GetExpressions(1)
	response := ExpressionsResponse{Expressions: expressions}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("cant marsahl response :(")
		w.WriteHeader(500)
		return
	}
	w.Write(jsonBytes)
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}
type CalculateResponse struct {
	Id int64 `json:"id"`
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
		log.Println("cant read body :(")
		w.WriteHeader(500)
	}
	defer r.Body.Close()
	var request CalculateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Println("cant unmarsahl body :(")
		w.WriteHeader(500)
	}
	if !isValidExpression(request.Expression) {
		w.WriteHeader(422)
		return
	}
	// Брать userID из заголовка
	id := logic.NewEx(request.Expression, r.Header.Get("userID"))
	response := CalculateResponse{Id: id}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("cant marsahl response :(")
		w.WriteHeader(500)
	}
	w.WriteHeader(201)
	w.Write(jsonBytes)
}

type GetExpressionResponse struct {
	Expression sqlite.Expression
}

func GetExpressionByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	userID := r.Header.Get("userID")
	userID1, _ := strconv.Atoi(userID)
	ex := storage.DataBase.GetExpressionById(int64(id), int64(userID1))
	if err != nil {
		w.WriteHeader(404)
	}
	response := GetExpressionResponse{ex}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("cant marsahl response :(")
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
			log.Println("cant read body :(")
			w.WriteHeader(500)
			return
		}
		defer r.Body.Close()
		var request GetSetTaskRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Println("cant unmarsahl body :(")
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

type RegisterLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("cant read body :(")
			w.WriteHeader(500)
			return
		}

		defer r.Body.Close()
		var request RegisterLoginRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Println("cant unmarsahl body, %v", err)
			w.WriteHeader(500)
		}
		login := request.Login
		password := request.Password
		// Проверка на существование пользователя с данным логином
		if ok, _ := storage.DataBase.UserExists(login); ok {
			log.Println("юзер с данным логином уже существует")
			w.WriteHeader(500)
		}
		hashedPass, err := logic.Generate(password)
		if err != nil {
			log.Println("cant hashing pass")
			w.WriteHeader(500)
		}
		id, err := storage.DataBase.AddUser(login, hashedPass)
		if err != nil {
			log.Println("cant add user")
			w.WriteHeader(500)
		}
		//fmt.Println(id)
		r.Header.Set("id", string(id))
		//cookieToken, err := r.Cookie("token")
		//a := cookieToken.Value
	}
}

type User struct {
	Login    string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		var u User
		json.NewDecoder(r.Body).Decode(&u)
		login := u.Login
		password := u.Password
		if ok, _ := storage.DataBase.UserExists(login); !ok {
			log.Println("Пользователя с данным логином не существует")
			w.WriteHeader(500)
			w.Write([]byte("Пользователя с данным логином не существует"))
			return
		}
		userFromDB, err := storage.DataBase.GetUser(login)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Пользователя с данным логином не существует"))
			return
		}
		err = logic.ComparePassword(userFromDB.Password, password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Неверный пароль")
			w.Write([]byte("Неверный пароль"))
			return
		}
		tokenString, err := logic.CreateToken(userFromDB.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("No username found")
			w.Write([]byte("No username found"))
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokenString)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
