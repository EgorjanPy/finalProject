package server_test

// import (
// 	"finalProject/pkg/calculation"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"testing"

// 	"net/http/httptest"
// )

// type TestingRequest struct {
// 	id         int
// 	expression string
// 	expected   string
// 	statusCode int
// }

// func TestCalcHandlerSuccessCase(t *testing.T) {
// 	requests := []TestingRequest{
// 		{1, "2+2*2", fmt.Sprintf("result: %f", 6.), 200},
// 		{2, "2+2", fmt.Sprintf("result: %f", 4.), 200},
// 		{3, "2+2*(2/2)", fmt.Sprintf("result: %f", 4.), 200},
// 		{4, "2+0", fmt.Sprintf("result: %f", 2.), 200},
// 	}
// 	for _, r := range requests {
// 		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", nil)
// 		w := httptest.NewRecorder()
// 		w.Header().Set("expression", r.expression)
// 		CalcHandler(w, req)
// 		res := w.Result()
// 		defer res.Body.Close()
// 		data, err := io.ReadAll(res.Body)
// 		if string(data) != r.expected {
// 			t.Errorf("Error: %v, expected %s, but got %s", err, r.expected, string(data))
// 		}
// 		if res.StatusCode != r.statusCode {
// 			t.Errorf("wrong status code, expected %d, but got %d", res.StatusCode, r.statusCode)
// 		}

// 	}
// }

// func TestCalcHandlerBadRequestCase(t *testing.T) {
// 	requests := []TestingRequest{
// 		{1, "2+2*(2", fmt.Sprintf("%v", calculation.ErrInvalidExpression), 422},
// 		{2, "2+2)", calculation.ErrInvalidExpression.Error(), 422},
// 		{3, "2+2*(2/0)", calculation.ErrDivisionByZero.Error(), 422},
// 		{4, "", calculation.ErrEmptyExpression.Error(), 422},
// 		{5, "*2+0", calculation.ErrInvalidExpression.Error(), 422},
// 		{6, "2+0*a", calculation.ErrInvalidExpression.Error(), 422},
// 		{7, "2+(5*3-+)", calculation.ErrInvalidExpression.Error(), 422},
// 	}
// 	for _, r := range requests {
// 		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", nil)
// 		w := httptest.NewRecorder()
// 		w.Header().Set("expression", r.expression)
// 		CalcHandler(w, req)
// 		res := w.Result()
// 		defer res.Body.Close()
// 		data, _ := io.ReadAll(res.Body)
// 		if string(data) != r.expected {
// 			t.Errorf("id %d wrong result: expected %s, but got %s", r.id, r.expected, string(data))
// 		}
// 		if res.StatusCode != r.statusCode {
// 			t.Errorf("id %d wrong status code, expected %d, but got %d", r.id, r.statusCode, res.StatusCode)
// 		}

// 	}
// }
// func TestCalcHandlerBadMethodCase(t *testing.T) {
// 	requests := []TestingRequest{
// 		{1, "2+2", "Method Not Allowed", 405},
// 		{2, "2+2*2", "Method Not Allowed", 405},
// 		{3, "2+2/0", "Method Not Allowed", 405},
// 		{4, "", "Method Not Allowed", 405},
// 		{5, "2+0", "Method Not Allowed", 405},
// 		{6, "2+0*1", "Method Not Allowed", 405},
// 		{7, "2+(5*3)", "Method Not Allowed", 405},
// 	}
// 	for _, r := range requests {
// 		req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
// 		w := httptest.NewRecorder()
// 		w.Header().Set("expression", r.expression)
// 		CalcHandler(w, req)
// 		res := w.Result()
// 		defer res.Body.Close()
// 		data, _ := io.ReadAll(res.Body)
// 		if string(data) != r.expected {
// 			t.Errorf("id %d wrong result: expected %s, but got %s", r.id, r.expected, string(data))
// 		}
// 		if res.StatusCode != r.statusCode {
// 			t.Errorf("id %d wrong status code, expected %d, but got %d", r.id, r.statusCode, res.StatusCode)
// 		}

// 	}
// }

// /// invalid expression
// /// invalid expression
