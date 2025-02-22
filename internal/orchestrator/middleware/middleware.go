package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
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
		if r.Method != http.MethodPost {
			// http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
			// return
			// TODO
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			// TODO
		}
		defer r.Body.Close()

		var request Request
		err = json.Unmarshal(body, &request)
		if err != nil || request.Expression == "" {
			w.WriteHeader(422)
			return
		}
		if !Validate(request.Expression) {
			w.WriteHeader(422)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func CalcLogger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		// request := new(Request)
		// defer r.Body.Close()
		// err := json.NewDecoder(r.Body).Decode(&request)
		// if err != nil {
		// 	logger.Error("finished",
		// 		slog.Group("req",
		// 			slog.String("method", r.Method),
		// 			slog.String("url", r.URL.String()),
		// 			slog.String("expression", request.Expression),
		// 		),
		// 		slog.Int("status", 500),
		// 		slog.Duration("duration", time.Second))
		// 	w.WriteHeader(500)
		// 	fmt.Fprintf(w, "Internal server error")
		// 	return
		// }
		// logger.Info("finished",
		// 	slog.Group("req",
		// 		slog.String("method", r.Method),
		// 		slog.String("url", r.URL.String()),
		// 		slog.String("expression", request.Expression)),
		// 	slog.Int("status", http.StatusOK),
		// 	slog.Duration("duration", time.Second))
		// w.Header().Set("expression", request.Expression)
		// next.ServeHTTP(w, r)
	})
}
