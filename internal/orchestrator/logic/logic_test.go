package logic

import (
	"finalProject/internal/storage/sqlite"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage реализует интерфейс, который используется в вашем коде
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SetResult(id int64, res string) error {
	args := m.Called(id, res)
	return args.Error(0)
}

// Добавляем другие необходимые методы хранилища
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

func TestExpressionEvaluation(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"Simple addition", "2+3", 5},
		{"Simple subtraction", "5-2", 3},
		{"Simple multiplication", "2*3", 6},
		{"Simple division", "6/2", 3},
		{"Complex expression", "2+3*4", 14},
		{"With parentheses", "(2+3)*4", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.expression)
			ast := parser.ParseExpression()
			result := ast.Evaluate()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSaveTasks(t *testing.T) {
	tasks := SaveTasks{Tasks: make(map[int]Task)}

	t.Run("Add and Get Task", func(t *testing.T) {
		testTask := Task{Arg1: 2, Arg2: 3, Operation: "+"}
		tasks.AddTask(1, testTask)

		retrievedTask, err := tasks.GetTaskById(1)
		assert.NoError(t, err)
		assert.Equal(t, testTask, retrievedTask)
	})

	t.Run("Get non-existent Task", func(t *testing.T) {
		_, err := tasks.GetTaskById(999)
		assert.Error(t, err)
	})
}

func TestSaveResults(t *testing.T) {
	results := SaveResults{Results: make(map[int]float64)}

	t.Run("Set and Get Result", func(t *testing.T) {
		results.SetResult(1, 42.0)
		assert.True(t, results.IsExists(1))
		assert.Equal(t, 42.0, results.GetResult(1))
	})

	t.Run("Get non-existent Result", func(t *testing.T) {
		assert.False(t, results.IsExists(999))
		assert.Equal(t, 0.0, results.GetResult(999))
	})
}

func TestPasswordHashing(t *testing.T) {
	password := "securepassword123"

	t.Run("GeneratePasswordHash and Compare", func(t *testing.T) {
		hash, err := GeneratePasswordHash(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		err = ComparePassword(hash, password)
		assert.NoError(t, err)
	})

	t.Run("Compare wrong password", func(t *testing.T) {
		hash, _ := GeneratePasswordHash(password)
		err := ComparePassword(hash, "wrongpassword")
		assert.Error(t, err)
	})
}

func TestJWT(t *testing.T) {
	userID := "123"

	t.Run("Create and Verify Token", func(t *testing.T) {
		token, err := CreateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		err = VerifyToken(token)
		assert.NoError(t, err)
	})

	t.Run("Get Payload from Token", func(t *testing.T) {
		token, _ := CreateToken(userID)
		payload, ok := JwtPayloadsFromToken(token)
		assert.True(t, ok)
		assert.Equal(t, userID, payload["sub"])
	})

	t.Run("Verify Invalid Token", func(t *testing.T) {
		err := VerifyToken("invalid.token.here")
		assert.Error(t, err)
	})
}
