package handlers

import (
	"encoding/json"
	"errors"
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
	"unicode"

	"github.com/gorilla/mux"
)

func hasDivisionByZero(expression string) bool {
	re := regexp.MustCompile(`\/0(\.0*)?([^0-9]|$)`)
	return re.MatchString(expression)
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
	return !hasDivisionByZero(cleanedExpression)
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
func GetUserID(r *http.Request) (string, error) {
	token, err := r.Cookie("jwtToken")
	if err != nil {
		return "", errors.New("error: cookie not found")
	}
	tokenString := token.String()[9:]
	jwtPayload, ok := logic.JwtPayloadsFromToken(tokenString)
	if !ok {
		return "", errors.New("error: invalid token claims")
	}
	userID, ok := jwtPayload["sub"].(string)
	if !ok {
		return "", errors.New("error: cant find sub from claims")
	}
	return userID, nil
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
	Id         int64  `json:"id"`
	Message    string `json:"message"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CalculateHandlers"
	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp := CalculateResponse{StatusCode: http.StatusBadRequest, Id: 0, Message: "bad request"}
		Response(w, resp)
		log.Fatalf("%s error: %v", op, err)
	}
	defer r.Body.Close()

	var request CalculateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		resp := CalculateResponse{StatusCode: http.StatusBadRequest, Id: 0, Message: "bad request"}
		Response(w, resp)
		return
	}
	if !isValidExpression(request.Expression) {
		resp := CalculateResponse{StatusCode: http.StatusBadRequest, Id: 0, Message: "wrong expression"}
		Response(w, resp)
		return
	}
	userID, err := GetUserID(r)
	if err != nil {
		resp := CalculateResponse{StatusCode: http.StatusBadRequest, Id: 0, Message: "wrong jwt token"}
		Response(w, resp)
		return
	}
	id, err := storage.DataBase.AddExpression(&sqlite.Expression{UserID: userID, Expression: request.Expression})
	if err != nil {
		resp := CalculateResponse{http.StatusInternalServerError, id, "cant add expression"}
		Response(w, resp)
		log.Fatalf("%s error: %v", op, err)
	}
	logic.NewExpression(id, request.Expression, userID)
	resp := CalculateResponse{http.StatusOK, id, fmt.Sprintf("Expression %d was added", id)}
	Response(w, resp)
	return
}

type ExpressionsResponse struct {
	StatusCode  uint                `json:"status_code"`
	Expressions []sqlite.Expression `json:"expressions"`
	Message     string              `json:"message"`
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	op := "handlers.ExpressionsHandler"
	userID, err := GetUserID(r)
	if err != nil {
		resp := ExpressionsResponse{StatusCode: http.StatusUnauthorized, Expressions: []sqlite.Expression{}, Message: "wrong jwt token"}
		Response(w, resp)
		return
	}
	intUserID, _ := strconv.Atoi(userID)
	expressions, err := storage.DataBase.GetExpressions(int64(intUserID))
	if err != nil {
		resp := ExpressionsResponse{StatusCode: http.StatusUnauthorized, Expressions: []sqlite.Expression{}, Message: "cant get expressions"}
		Response(w, resp)
		log.Fatalf("%s error: %v", op, err)
	}
	resp := ExpressionsResponse{StatusCode: http.StatusOK, Expressions: expressions, Message: "success"}
	Response(w, resp)
	return
}
func GetExpressionByIdHandler(w http.ResponseWriter, r *http.Request) {
	strId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		resp := ExpressionsResponse{StatusCode: http.StatusNotFound, Expressions: make([]sqlite.Expression, 0), Message: "wrong id"}
		Response(w, resp)
		return
	}
	userID, err := GetUserID(r)
	if err != nil {
		resp := ExpressionsResponse{StatusCode: http.StatusUnauthorized, Expressions: []sqlite.Expression{}, Message: "wrong jwt token"}
		Response(w, resp)
		return
	}
	intUserID, _ := strconv.Atoi(userID)
	ex, err := storage.DataBase.GetExpressionById(int64(id), int64(intUserID))
	if err != nil {
		resp := ExpressionsResponse{StatusCode: http.StatusNotFound, Expressions: make([]sqlite.Expression, 0), Message: "expression not found"}
		Response(w, resp)
		log.Println(err)
		return
	}
	resp := ExpressionsResponse{StatusCode: http.StatusOK, Expressions: []sqlite.Expression{ex}, Message: "success"}
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			resp := LoginRegisterResponse{StatusCode: http.StatusBadRequest, Message: "bad request"}
			Response(w, resp)
			return
		}
		defer r.Body.Close()
		var request RegisterLoginRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			resp := LoginRegisterResponse{StatusCode: http.StatusBadRequest, Message: "bad request"}
			Response(w, resp)
			return
		}
		login := request.Login
		password := request.Password

		if ok, _ := storage.DataBase.UserExists(login); ok {
			resp := LoginRegisterResponse{StatusCode: http.StatusUnauthorized, Message: "user already exists"}
			Response(w, resp)
			return
		}
		hashedPass, err := logic.Generate(password)
		if err != nil {
			resp := LoginRegisterResponse{StatusCode: http.StatusInternalServerError, Message: "cant generate jwt token"}
			Response(w, resp)
			log.Fatalf("%s error: %v", op, err)
		}
		_, err = storage.DataBase.AddUser(login, hashedPass)
		if err != nil {
			resp := LoginRegisterResponse{http.StatusInternalServerError, "cant add user"}
			Response(w, resp)
			log.Fatalf("%s error: %v", op, err)
		}
		resp := LoginRegisterResponse{http.StatusOK, "success"}
		Response(w, resp)
		return
	} else {
		resp := LoginRegisterResponse{http.StatusMethodNotAllowed, "method not allowed"}
		Response(w, resp)
		return
	}
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
			resp := LoginRegisterResponse{http.StatusUnauthorized, "user not found"}
			Response(w, resp)
			return
		}
		userFromDB, err := storage.DataBase.GetUser(login)
		if err != nil {
			resp := LoginRegisterResponse{http.StatusInternalServerError, "cant get user"}
			Response(w, resp)
			log.Fatalf("%s Error: %v", op, err)
		}
		err = logic.ComparePassword(userFromDB.Password, password)
		if err != nil {
			resp := LoginRegisterResponse{http.StatusUnauthorized, "wrong password"}
			Response(w, resp)
			return
		}
		tokenString, err := logic.CreateToken(userFromDB.ID)
		if err != nil {
			resp := LoginRegisterResponse{http.StatusInternalServerError, "cant create jwt token"}
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
		resp := LoginRegisterResponse{http.StatusOK, "auth success"}
		Response(w, resp)
		return
	} else {
		resp := LoginRegisterResponse{http.StatusMethodNotAllowed, "method not allowed"}
		Response(w, resp)
		return
	}
}
