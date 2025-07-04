package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		fmt.Printf("cant open db, %v", err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("cant ping db, %v", err)
		os.Exit(1)
	}
	err = CreateTables(db)
	if err != nil {
		fmt.Printf("cant create tables, %v", err)
		os.Exit(1)
	}
	return &Storage{db}, nil
}
func CreateTables(db *sql.DB) error {
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		login TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		expression TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		answer TEXT,
		status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'calculating', 'completed', 'failed')),
		FOREIGN KEY (user_id)  REFERENCES user(id) ON DELETE CASCADE
	);`
	)
	if _, err := db.Exec(usersTable); err != nil {
		return err
	}
	if _, err := db.Exec(expressionsTable); err != nil {
		return err
	}
	return nil
}

type (
	User struct {
		ID       string
		Login    string
		Password string
	}

	Expression struct {
		ID         int64
		Expression string
		Answer     sql.NullString
		Status     string
		UserID     string
	}
)

func (s *Storage) UserExists(login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`
	err := s.db.QueryRow(query, login).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
func (s *Storage) AddUser(login, password string) (int64, error) {
	var q = `
	INSERT INTO users (login, password) values ($1, $2)
	`
	result, err := s.db.Exec(q, login, password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) AddExpression(expression *Expression) (int64, error) {
	var q = `
	INSERT INTO expressions (expression, user_id) values ($1, $2)
	`
	result, err := s.db.Exec(q, expression.Expression, expression.UserID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (s *Storage) GetExpressions(id int64) ([]Expression, error) {
	var expressions []Expression
	var q = `SELECT * FROM expressions WHERE user_id = $1`
	rows, err := s.db.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID, &e.Answer, &e.Status)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}
	return expressions, nil
}
func (s *Storage) GetExpressionById(ex_id int64, user_id int64) (Expression, error) {
	var q = "SELECT * FROM expressions WHERE id = $1 AND user_id = $2"
	ex := Expression{}
	err := s.db.QueryRow(q, ex_id, user_id).Scan(&ex.ID, &ex.Expression, &ex.UserID, &ex.Answer, &ex.Status)
	if err != nil {
		return Expression{}, err
	}
	return ex, nil
}
func (s *Storage) GetUncompletedExpressions() ([]Expression, error) {
	expressions := []Expression{}
	var q = "SELECT * FROM expressions WHERE answer IS NULL"
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID, &e.Answer, &e.Status)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}
	return expressions, nil
}

func (s *Storage) UpdateUserPassword(id int64, pass string) error {
	var q = "UPDATE users SET password = $1 WHERE id = $2"
	_, err := s.db.Exec(q, pass, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) SetResult(id int64, res string) error {
	var q = `
        UPDATE expressions 
        SET answer = $1, 
            status = 'completed'
        WHERE id = $2
    `
	_, err := s.db.Exec(q, res, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(login string) (User, error) {
	var (
		user User
		err  error
	)
	var q = "SELECT id, login, password FROM users WHERE login=$1"
	err = s.db.QueryRow(q, login).Scan(&user.ID, &user.Login, &user.Password)
	return user, err
}
