package handlers

import (
	"encoding/json"
	"finalProject/internal/config"
	"finalProject/internal/orchestrator/logic"
	"finalProject/internal/storage"
	"finalProject/internal/storage/sqlite"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
)

var cfg = config.MustLoad()

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
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
func GetUserID(r *http.Request) (string, bool) {
	token, err := r.Cookie("jwtToken")
	if err != nil {
		log.Printf("autorixation error %v", err)
	}
	tokenString := token.String()[9:]
	jwtPayload, ok := logic.JwtPayloadsFromToken(tokenString)
	if !ok {
		log.Println("invalid token claims")
		return "", false
	}
	userID, ok := jwtPayload["sub"].(string)
	if !ok {
		log.Println("cant find sub from claims")
		return "", false
	}
	return userID, true
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}
type CalculateResponse struct {
	Id int64 `json:"id"`
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
	userID, ok := GetUserID(r)
	if !ok {
		log.Println("cant get userID from claims")
		w.WriteHeader(401)
	}
	id, err := storage.DataBase.AddExpression(&sqlite.Expression{UserID: userID, Expression: request.Expression})
	logic.NewExpression(id, request.Expression, userID)

	response := CalculateResponse{Id: id}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("cant marsahl response :(")
		w.WriteHeader(500)
	}
	w.WriteHeader(201)
	w.Write(jsonBytes)
}

type ExpressionsResponse struct {
	Expressions []sqlite.Expression
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r)
	if !ok {
		log.Println("cant get userID from claims")
		w.WriteHeader(401)
	}
	intUserID, _ := strconv.Atoi(userID)
	expressions, err := storage.DataBase.GetExpressions(int64(intUserID))
	if err != nil {
		log.Printf("error %v", err)

	}
	response := ExpressionsResponse{Expressions: expressions}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("cant marsahl response :(")
		w.WriteHeader(500)
		return
	}
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
	userID, ok := GetUserID(r)
	if !ok {
		log.Println("cant get userID from claims")
		w.WriteHeader(401)
	}
	intUserID, _ := strconv.Atoi(userID)
	ex, err := storage.DataBase.GetExpressionById(int64(id), int64(intUserID))
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

type RegisterLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	op := "handlers.RegisterHandler"
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			Response(w, http.StatusBadRequest, "bad request")
			return
		}
		defer r.Body.Close()
		var request RegisterLoginRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			Response(w, http.StatusBadRequest, "bad request")
			return
		}
		login := request.Login
		password := request.Password

		if ok, _ := storage.DataBase.UserExists(login); ok {
			Response(w, http.StatusUnauthorized, "user already exists")
			return
		}
		hashedPass, err := logic.Generate(password)
		if err != nil {
			Response(w, http.StatusInternalServerError, "cant generate jwt token")
			log.Fatalf("%s error: %v", op, err)
		}
		_, err = storage.DataBase.AddUser(login, hashedPass)
		if err != nil {
			Response(w, http.StatusInternalServerError, "cant add user")
			log.Fatalf("%s error: %v", op, err)
		}
	} else {
		Response(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

type LoginRegisterResponse struct {
	StatusCode uint
	Message    string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		op := "handlers.LoginHandler"
		w.Header().Set("Content-Type", "application/json")
		var u RegisterLoginRequest
		json.NewDecoder(r.Body).Decode(&u)
		login := u.Login
		password := u.Password
		if ok, _ := storage.DataBase.UserExists(login); !ok {
			Response(w, http.StatusUnauthorized, "user not found")
			return
		}
		userFromDB, err := storage.DataBase.GetUser(login)
		if err != nil {
			w.WriteHeader(500)
			log.Fatalf("%s Error: %v", op, err)
		}
		err = logic.ComparePassword(userFromDB.Password, password)
		if err != nil {
			Response(w, http.StatusUnauthorized, "wrong password")
			return
		}
		tokenString, err := logic.CreateToken(userFromDB.ID)
		if err != nil {
			w.WriteHeader(500)
			log.Fatalf("%s error: %v", op, err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "jwtToken",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,                    // Защита от доступа через JavaScript
			Secure:   true,                    // Отправка только по HTTPS
			SameSite: http.SameSiteStrictMode, // Защита от CSRF атак
		})
		Response(w, http.StatusOK, "auth success")
		return
	} else {
		Response(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}
func Response(w http.ResponseWriter, status uint, text string) {
	resp := LoginRegisterResponse{StatusCode: status, Message: text}
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("error: %v", err)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("error: %v", err)
	}
}
