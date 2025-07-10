package handlers

import (
	"bytes"
	"encoding/json"
	"finalProject/internal/storage"
	"finalProject/internal/storage/sqlite"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockStorage заменяет реальное хранилище для тестов
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UserExists(login string) (bool, error) {
	args := m.Called(login)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) AddUser(login, password string) (int64, error) {
	args := m.Called(login, password)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) GetUser(login string) (sqlite.User, error) {
	args := m.Called(login)
	return args.Get(0).(sqlite.User), args.Error(1)
}

func (m *MockStorage) AddExpression(expression *sqlite.Expression) (int64, error) {
	args := m.Called(expression)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) GetExpressions(id int64) ([]sqlite.Expression, error) {
	args := m.Called(id)
	return args.Get(0).([]sqlite.Expression), args.Error(1)
}

func (m *MockStorage) GetExpressionById(ex_id, user_id int64) (sqlite.Expression, error) {
	args := m.Called(ex_id, user_id)
	return args.Get(0).(sqlite.Expression), args.Error(1)
}

func TestHandlers(t *testing.T) {
	// Создаем мок хранилища
	mockStorage := new(MockStorage)
	storage.DataBase = mockStorage

	t.Run("Test isValidExpression", func(t *testing.T) {
		assert.True(t, isValidExpression("2+2"))
		assert.True(t, isValidExpression("(2+3)*4"))
		assert.False(t, isValidExpression("2/0"))
		assert.False(t, isValidExpression("2++2"))
		assert.False(t, isValidExpression("(2+3"))
	})

	t.Run("Test CalculateHandler with valid expression", func(t *testing.T) {
		// Настраиваем мок
		mockStorage.On("AddExpression", mock.Anything).Return(int64(1), nil)

		reqBody := CalculateRequest{Expression: "2+2"}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
		req.AddCookie(&http.Cookie{Name: "jwtToken", Value: "token=valid"})
		w := httptest.NewRecorder()

		CalculateHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp CalculateResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, int64(1), resp.Id)
	})

	t.Run("Test CalculateHandler with invalid expression", func(t *testing.T) {
		reqBody := CalculateRequest{Expression: "2++2"}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
		w := httptest.NewRecorder()

		CalculateHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test ExpressionsHandler", func(t *testing.T) {
		// Настраиваем мок
		expected := []sqlite.Expression{
			{ID: 1, Expression: "2+2", Status: "completed"},
		}
		mockStorage.On("GetExpressions", int64(1)).Return(expected, nil)

		req := httptest.NewRequest("GET", "/expressions", nil)
		req.AddCookie(&http.Cookie{Name: "jwtToken", Value: "token=valid"})
		w := httptest.NewRecorder()

		ExpressionsHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ExpressionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, len(resp.Expressions))
	})

	t.Run("Test GetExpressionByIdHandler", func(t *testing.T) {
		// Настраиваем мок
		expected := sqlite.Expression{ID: 1, Expression: "2+2", Status: "completed"}
		mockStorage.On("GetExpressionById", int64(1), int64(1)).Return(expected, nil)

		req := httptest.NewRequest("GET", "/expressions/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.AddCookie(&http.Cookie{Name: "jwtToken", Value: "token=valid"})
		w := httptest.NewRecorder()

		GetExpressionByIdHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ExpressionsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, len(resp.Expressions))
	})

	t.Run("Test RegisterHandler success", func(t *testing.T) {
		// Настраиваем мок
		mockStorage.On("UserExists", "newuser").Return(false, nil)
		mockStorage.On("AddUser", "newuser", mock.Anything).Return(int64(1), nil)

		reqBody := RegisterLoginRequest{Login: "newuser", Password: "pass"}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		w := httptest.NewRecorder()

		RegisterHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Test LoginHandler success", func(t *testing.T) {
		// Настраиваем мок
		user := sqlite.User{ID: "1", Login: "user", Password: "hashedpass"}
		mockStorage.On("UserExists", "user").Return(true, nil)
		mockStorage.On("GetUser", "user").Return(user, nil)

		reqBody := RegisterLoginRequest{Login: "user", Password: "pass"}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		LoginHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	})
}
