package handlers

import (
	"encoding/json"
	"errors"
	"finalProject/internal/orchestrator/logic"
	"finalProject/internal/orchestrator/tools"
	"finalProject/internal/storage"
	"finalProject/internal/storage/sqlite"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
)

func hasDivisionByZero(expression string) error {
	re := regexp.MustCompile(`\/0(\.0*)?([^0-9]|$)`)
	if re.MatchString(expression) {
		return errors.New("division by zero")
	}
	return nil
}
func isValidExpression(expression string) error {
	cleanedExpression := removeSpaces(expression)
	validPattern := `^[0-9+\-*/()]+$`
	matched, err := regexp.MatchString(validPattern, cleanedExpression)
	if err != nil || !matched {
		return errors.New("invalid expression")
	}
	if !areParenthesesBalanced(cleanedExpression) {
		return errors.New("wrong parentheses")
	}
	if strings.Contains("+=-/*:", string(expression[0])) || strings.Contains("+=-/*:", string(expression[len(expression)-1])) {
		return errors.New("invalid expression")
	}
	if err := hasDivisionByZero(cleanedExpression); err != nil {
		return err
	}
	return nil
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

func Response(w http.ResponseWriter, resp any) {
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

type CalculateRequest struct {
	Expression string `json:"expression"`
}
type CalculateResponse struct {
	StatusCode uint   `json:"status_code"`
	Id         int    `json:"id"`
	Message    string `json:"message"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CalculateHandlers"
	userID, err := tools.GetUserIDFromContext(r)
	if err != nil {
		message := err.Error()
		resp := CalculateResponse{StatusCode: http.StatusUnauthorized, Id: 0, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, resp)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		message := "bad request: can't read request body"
		response := CalculateResponse{StatusCode: http.StatusBadRequest, Id: 0, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, response)
		return
	}
	defer r.Body.Close()

	var request CalculateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		message := "bad request: can't unmarshal body"
		response := CalculateResponse{StatusCode: http.StatusBadRequest, Id: 0, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, response)
		return
	}
	if err := isValidExpression(request.Expression); err != nil {
		message := err.Error()
		response := CalculateResponse{StatusCode: http.StatusBadRequest, Id: 0, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, response)
		return
	}
	expressionID, err := storage.DataBase.AddExpression(&sqlite.Expression{UserID: userID, Expression: request.Expression})
	if err != nil {
		response := CalculateResponse{http.StatusInternalServerError, expressionID, "cant add expression"}
		Response(w, response)
		log.Fatalf("%s error: %v", op, err)
	}
	logic.NewExpression(expressionID, request.Expression, userID)
	message := fmt.Sprintf("expression %d added", expressionID)
	response := CalculateResponse{http.StatusOK, expressionID, message}
	log.Printf("%s, %s", r.URL, message)
	Response(w, response)
	return
}

type ExpressionsResponse struct {
	StatusCode  uint                 `json:"status_code"`
	Expressions *[]sqlite.Expression `json:"expressions,omitempty"`
	Message     string               `json:"message"`
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	op := "handlers.ExpressionsHandler"
	userID, err := tools.GetUserIDFromContext(r)
	if err != nil {
		message := err.Error()
		response := ExpressionsResponse{StatusCode: http.StatusUnauthorized, Expressions: nil, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, response)
		return
	}
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		message := "wrong userID"
		response := ExpressionsResponse{StatusCode: http.StatusUnauthorized, Expressions: nil, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, response)
		return
	}
	expressions, err := storage.DataBase.GetExpressions(intUserID)
	if err != nil {
		message := "can't get expressions"
		resp := ExpressionsResponse{StatusCode: http.StatusBadRequest, Expressions: nil, Message: message}
		Response(w, resp)
		log.Fatalf("%s error: %v", op, err)
	}
	message := "success"
	response := ExpressionsResponse{StatusCode: http.StatusOK, Expressions: &expressions, Message: message}
	log.Printf("%s, %s", r.URL, message)
	Response(w, response)
	return
}
func GetExpressionByIdHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := tools.GetUserIDFromContext(r)
	if err != nil {
		message := err.Error()
		resp := ExpressionsResponse{StatusCode: http.StatusUnauthorized, Expressions: nil, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, resp)
		return
	}
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		message := "wrong userID"
		response := ExpressionsResponse{StatusCode: http.StatusUnauthorized, Expressions: nil, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, response)
		return
	}
	expressionID, err := strconv.Atoi(mux.Vars(r)["expressionID"])
	if err != nil {
		message := "wrong expressionID"
		resp := ExpressionsResponse{StatusCode: http.StatusNotFound, Expressions: nil, Message: message}
		log.Printf("%s, error: %s", r.URL, message)
		Response(w, resp)
		return
	}

	expression, err := storage.DataBase.GetExpressionById(expressionID, intUserID)
	if err != nil {
		message := "expression not found"
		resp := ExpressionsResponse{StatusCode: http.StatusNotFound, Expressions: nil, Message: message}
		log.Printf("%s, error: %s", r.URL, err)
		Response(w, resp)
		return
	}
	message := "success"
	resp := ExpressionsResponse{StatusCode: http.StatusOK, Expressions: &[]sqlite.Expression{expression}, Message: message}
	log.Printf("%s, %s", r.URL, message)
	Response(w, resp)
	return
}

type RegisterLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type LoginRegisterResponse struct {
	StatusCode uint   `json:"status_code"`
	Message    string `json:"message"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	op := "handlers.RegisterHandler"
	if r.Method == http.MethodPost {
		var request RegisterLoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			message := "bad request"
			response := LoginRegisterResponse{StatusCode: http.StatusBadRequest, Message: message}
			log.Printf("%s, error: %s", r.URL, message)
			Response(w, response)
			return
		}
		login := request.Login
		password := request.Password

		if ok, err := storage.DataBase.UserExists(login); ok || err != nil {
			message := "user exists"
			response := LoginRegisterResponse{StatusCode: http.StatusUnauthorized, Message: message}
			log.Printf("%s, error: %s", r.URL, message)
			Response(w, response)
			return
		}
		hashedPass, err := tools.GeneratePasswordHash(password)
		if err != nil {
			message := "can't generate jwt token"
			response := LoginRegisterResponse{StatusCode: http.StatusInternalServerError, Message: message}
			Response(w, response)
			log.Fatalf("%s error: %v", op, err)
		}
		_, err = storage.DataBase.AddUser(login, hashedPass)
		if err != nil {
			message := "cant add user"
			response := LoginRegisterResponse{http.StatusInternalServerError, message}
			log.Printf("%s, %s", r.URL, message)
			Response(w, response)
			log.Fatalf("%s error: %v", op, err)
		}
		message := "success"
		response := LoginRegisterResponse{http.StatusOK, message}
		log.Printf("%s, %s", r.URL, message)
		Response(w, response)
		return
	} else {
		message := "method not allowed"
		response := LoginRegisterResponse{http.StatusMethodNotAllowed, message}
		log.Printf("%s, %s", r.URL, message)
		Response(w, response)
		return
	}
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		op := "handlers.LoginHandler"
		var request RegisterLoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			message := "bad request"
			response := LoginRegisterResponse{StatusCode: http.StatusBadRequest, Message: message}
			log.Printf("%s, error: %s", r.URL, message)
			Response(w, response)
			return
		}
		login := request.Login
		password := request.Password
		if ok, err := storage.DataBase.UserExists(login); !ok || err != nil {
			message := "user not found"
			resp := LoginRegisterResponse{http.StatusUnauthorized, message}
			log.Printf("%s, %s", r.URL, message)
			Response(w, resp)
			return
		}
		userFromDB, err := storage.DataBase.GetUser(login)
		if err != nil {
			message := "can't get user"
			resp := LoginRegisterResponse{http.StatusInternalServerError, message}
			Response(w, resp)
			log.Fatalf("%s Error: %v", op, err)
		}
		err = tools.ComparePassword(userFromDB.Password, password)
		if err != nil {
			message := "wrong password"
			resp := LoginRegisterResponse{http.StatusUnauthorized, message}
			log.Printf("%s, %s", r.URL, message)
			Response(w, resp)
			return
		}
		tokenString, err := tools.CreateToken(userFromDB.ID)
		if err != nil {
			message := "can't create token"
			resp := LoginRegisterResponse{http.StatusInternalServerError, message}
			Response(w, resp)
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
		message := "auth success"
		resp := LoginRegisterResponse{http.StatusOK, message}
		log.Printf("%s, %s", r.URL, message)
		Response(w, resp)
		return
	} else {
		message := "method not allowed"
		resp := LoginRegisterResponse{http.StatusMethodNotAllowed, message}
		log.Printf("%s, %s", r.URL, message)
		Response(w, resp)
		return
	}
}
