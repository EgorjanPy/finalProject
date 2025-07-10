package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	suite.Suite
	storage *Storage
	dbPath  string
}

func (s *StorageTestSuite) SetupTest() {
	s.dbPath = "test.db"
	var err error
	s.storage, err = New(s.dbPath)
	assert.NoError(s.T(), err)
}

func (s *StorageTestSuite) TearDownTest() {
	s.storage.db.Close()
	os.Remove(s.dbPath)
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestUserExists() {
	// Test non-existing user
	exists, err := s.storage.UserExists("nonexistent")
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)

	// Add test user
	_, err = s.storage.AddUser("testuser", "password")
	assert.NoError(s.T(), err)

	// Test existing user
	exists, err = s.storage.UserExists("testuser")
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *StorageTestSuite) TestAddUser() {
	// Add new user
	id, err := s.storage.AddUser("newuser", "pass123")
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), id, int64(0))

	// Try to add duplicate user
	_, err = s.storage.AddUser("newuser", "pass123")
	assert.Error(s.T(), err)
}

func (s *StorageTestSuite) TestAddAndGetExpression() {
	// Add test user
	userID, err := s.storage.AddUser("expruser", "pass")
	assert.NoError(s.T(), err)

	// Add expression
	expr := &Expression{
		Expression: "2+2",
		UserID:     string(fmt.Sprint(userID)),
	}
	exprID, err := s.storage.AddExpression(expr)
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), exprID, int64(0))

	// Get expressions
	exprs, err := s.storage.GetExpressions(userID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), exprs, 1)
	assert.Equal(s.T(), "2+2", exprs[0].Expression)
	assert.Equal(s.T(), "pending", exprs[0].Status)
}

func (s *StorageTestSuite) TestGetExpressionById() {
	// Add test data
	userID, _ := s.storage.AddUser("getbyid", "pass")
	exprID, _ := s.storage.AddExpression(&Expression{
		Expression: "3*3",
		UserID:     fmt.Sprint(userID),
	})

	// Test get by id
	expr, err := s.storage.GetExpressionById(exprID, userID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "3*3", expr.Expression)

	// Test wrong id
	_, err = s.storage.GetExpressionById(999, userID)
	assert.Error(s.T(), err)
}

func (s *StorageTestSuite) TestGetUncompletedExpressions() {
	// Add test data
	userID, _ := s.storage.AddUser("uncompleted", "pass")
	s.storage.AddExpression(&Expression{
		Expression: "1+1",
		UserID:     fmt.Sprint(userID),
	})
	s.storage.AddExpression(&Expression{
		Expression: "2+2",
		UserID:     fmt.Sprint(userID),
	})

	// Mark one as completed
	s.storage.SetResult(1, "2")

	// Get uncompleted
	exprs, err := s.storage.GetUncompletedExpressions()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), exprs, 1)
	assert.Equal(s.T(), "2+2", exprs[0].Expression)
}

func (s *StorageTestSuite) TestUpdateUserPassword() {
	// Add test user
	userID, _ := s.storage.AddUser("changepass", "oldpass")

	// Update password
	err := s.storage.UpdateUserPassword(userID, "newpass")
	assert.NoError(s.T(), err)

	// Verify
	user, err := s.storage.GetUser("changepass")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "newpass", user.Password)
}

func (s *StorageTestSuite) TestSetResult() {
	// Add test data
	userID, _ := s.storage.AddUser("setresult", "pass")
	exprID, _ := s.storage.AddExpression(&Expression{
		Expression: "5+5",
		UserID:     fmt.Sprint(userID),
	})

	// Set result
	err := s.storage.SetResult(exprID, "10")
	assert.NoError(s.T(), err)

	// Verify
	expr, err := s.storage.GetExpressionById(exprID, userID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "10", expr.Answer.String)
	assert.Equal(s.T(), "completed", expr.Status)
}

func (s *StorageTestSuite) TestGetUser() {
	// Add test user
	_, err := s.storage.AddUser("getuser", "userpass")
	assert.NoError(s.T(), err)

	// Get user
	user, err := s.storage.GetUser("getuser")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "getuser", user.Login)
	assert.Equal(s.T(), "userpass", user.Password)

	// Test non-existing user
	_, err = s.storage.GetUser("notexists")
	assert.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, sql.ErrNoRows))
}
