package sqlite

import (
	"context"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func CreateTables(ctx context.Context, db *sql.DB) error {
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
	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}
	return nil
}

type (
	User struct {
		ID       int64
		Login    string
		Password string
	}

	Expression struct {
		ID         int64
		Expression string
		Answer     string
		Status     string
		UserID     int64
	}
)

func AddUser(ctx context.Context, db *sql.DB, user *User) (int64, error) {
	var q = `
	INSERT INTO users (login, password) values ($1, $2)
	`
	result, err := db.ExecContext(ctx, q, user.Login, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func AddExpression(ctx context.Context, db *sql.DB, expression *Expression) (int64, error) {
	var q = `
	INSERT INTO expressions (expression, user_id) values ($1, $2)
	`
	result, err := db.ExecContext(ctx, q, expression.Expression, expression.UserID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
func GetExpressions(ctx context.Context, db *sql.DB, id int64) ([]Expression, error) {
	var expressions []Expression
	var q = "SELECT id, expression, user_id FROM expressions WHERE user_id = $1"

	rows, err := db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := Expression{}
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}
	return expressions, nil
}
func GetExpressionById(ctx context.Context, db *sql.DB, ex_id int64, user_id int64) Expression {
	var q = "SELECT expression, answer, status WHERE id = $1 AND user_id = $2"
	var ex Expression
	_ = db.QueryRowContext(ctx, q, ex_id, user_id).Scan(&ex.Expression, &ex.Answer, &ex.Status)
	return ex
}
func UpdateUserPassword(ctx context.Context, db *sql.DB, id int64, pass string) error {
	var q = "UPDATE users SET password = $1 WHERE id = $2"
	_, err := db.ExecContext(ctx, q, pass, id)
	if err != nil {
		return err
	}
	return nil
}
func SetResult(ctx context.Context, db *sql.DB, id int64, res string) error {
	var q = "UPDATE expressions SET answer = $1 WHERE id = $2"
	_, err := db.ExecContext(ctx, q, res, id)
	if err != nil {
		return err
	}
	return nil
}

// Insert user +
// Update password +

// Insert expression +
// Get expressions +
// Все функции из logic ?!
