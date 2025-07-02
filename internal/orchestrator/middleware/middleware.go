package middleware

import (
	"finalProject/internal/orchestrator/logic"
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

var secretKey = []byte("secret-key")

func ProtectedHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token, err := r.Cookie("jwtToken")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("autorixation error %v", err)

		}
		tokenString := token.String()[9:]
		err = logic.VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Invalid token")
			return
		}
		//fmt.Fprint(w, "Welcome to the the protected area")
		next.ServeHTTP(w, r)
	})
}
