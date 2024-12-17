package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"finalProject/pkg/calculation"
)

type Config struct {
	Port string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Функция запуска приложения
// тут будем читать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression or output to exit")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		// убираем пробелы, чтобы оставить только вычислемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		//вычисляем выражение
		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}

}

type Request struct {
	Expression string `json:"expression"`
}

func CalcLogger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		request := new(Request)
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.Error("finished",
				slog.Group("req",
					slog.String("method", r.Method),
					slog.String("url", r.URL.String()),
					slog.String("msg", "can't decode request"),
				),
				slog.Int("status", http.StatusBadRequest),
				slog.Duration("duration", time.Second))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Info("finished",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("expression", request.Expression)),
			slog.Int("status", http.StatusOK),
			slog.Duration("duration", time.Second))
		w.Header().Set("expression", request.Expression)
		next.ServeHTTP(w, r)
	})
}
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	expression := w.Header().Get("expression")
	res, err := calculation.Calc(expression)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%v", err)
		return
	}
	// log.Println(res)

	fmt.Fprintf(w, "result: %f", res)
}
func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcLogger(CalcHandler))
	// fmt.Printf("The server has started! Go to https://127.0.0.1:%s to view", a.config.Port)
	return http.ListenAndServe(":"+a.config.Port, nil)

}
