package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
		next.ServeHTTP(w, r)
	})
}
func CheckToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логику проверки токена
	})
}
