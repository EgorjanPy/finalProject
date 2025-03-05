package middleware

import (
	"encoding/json"
	"finalProject/internal/orchestrator/handlers"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Expression string `json:"expression"`
}

func Validate(expression string) bool {
	prohibitedSymbols := "qwertyuiopasdfghjklzxcvbnm!№;%?_#$^&~"
	for _, r := range expression {
		if strings.Contains(prohibitedSymbols, strings.ToLower(string(r))) {
			return false
		}
	}
	if strings.Contains("+=-/*:", string(expression[0])) || strings.Contains("+=-/*:", string(expression[len(expression)-1])) {
		return false
	}
	return true
}

func CalculateValidation(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
		}
		defer r.Body.Close()

		var request handlers.CalculateRequest
		err = json.Unmarshal(body, &request)
		if err != nil || request.Expression == "" {
			w.WriteHeader(500)
		}
		if !Validate(request.Expression) {
			w.WriteHeader(422)
		}
		fmt.Println("Всё гуд")
		next.ServeHTTP(w, r)
	})
}
func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
		next.ServeHTTP(w, r)
	})
}
