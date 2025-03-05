package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockServer создает тестовый сервер для эмуляции оркестратора
func MockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			response := Response{
				Id:             1,
				Arg1:           2,
				Arg2:           3,
				Operation:      "+",
				Operation_time: 100 * time.Millisecond,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if r.Method == http.MethodPost {
			var request Request
			json.NewDecoder(r.Body).Decode(&request)
			if request.Id != 1 || request.Result != 5 {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}))
}

func TestStartAgent(t *testing.T) {
	server := MockServer()
	defer server.Close()

	// Создаем агента
	app := New("", 1)
	app.Port = server.URL
	go app.StartAgent()

	time.Sleep(500 * time.Millisecond)
}

func TestStartApp(t *testing.T) {
	server := MockServer()
	defer server.Close()
	app := New("", 3)
	app.Port = server.URL
	go app.StartApp()

	time.Sleep(1 * time.Second)
}

func TestOperations(t *testing.T) {
	tests := []struct {
		arg1      float64
		arg2      float64
		operation string
		expected  float64
	}{
		{2, 3, "+", 5},
		{5, 3, "-", 2},
		{2, 3, "*", 6},
		{6, 3, "/", 2},
	}

	for _, test := range tests {
		var res float64
		switch test.operation {
		case "+":
			res = test.arg1 + test.arg2
		case "-":
			res = test.arg1 - test.arg2
		case "*":
			res = test.arg1 * test.arg2
		case "/":
			res = test.arg1 / test.arg2
		}

		if res != test.expected {
			t.Errorf("Operation %s failed: expected %f, got %f", test.operation, test.expected, res)
		}
	}
}
