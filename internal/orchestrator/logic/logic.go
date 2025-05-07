package logic

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Expression struct {
	Id         int
	Expression string
	Status     string
	Result     float64
}

type SaveExpressions struct {
	mu          sync.Mutex
	Expressions []Expression
}

var Expressions = SaveExpressions{
	mu:          sync.Mutex{},
	Expressions: []Expression{},
}

func (se *SaveExpressions) GetExpressions() []Expression {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.Expressions
}

func (se *SaveExpressions) SetResult(id int, res float64) {
	se.mu.Lock()
	se.Expressions[id].Result = res
	se.Expressions[id].Status = "complited"
	se.mu.Unlock()
}
func (se *SaveExpressions) AddExpression(ex Expression) {
	se.mu.Lock()
	se.Expressions = append(se.Expressions, ex)
	se.mu.Unlock()
}
func (se *SaveExpressions) GetExpressionById(id int) (Expression, error) {
	se.mu.Lock()
	defer se.mu.Unlock()
	for _, ex := range se.Expressions {
		if ex.Id == id {
			return ex, nil
		}
	}
	return Expression{}, fmt.Errorf("not found")
}

type Task struct {
	Id        int
	Arg1      float64
	Arg2      float64
	Operation string
}
type SaveTasks struct {
	mu    sync.Mutex
	Tasks map[int]Task
}

func (st *SaveTasks) GetLen() int {
	st.mu.Lock()
	defer st.mu.Unlock()
	return len(st.Tasks)
}
func (st *SaveTasks) AddTask(id int, task Task) {
	st.mu.Lock()
	st.Tasks[id] = task
	st.mu.Unlock()
}
func (st *SaveTasks) GetTaskById(id int) (Task, error) {
	if st.GetLen() > Results.GetLen() {
		st.mu.Lock()
		defer st.mu.Unlock()
		if t, exists := st.Tasks[id]; exists {
			return t, nil
		}
		return Task{}, fmt.Errorf("not found")
	}
	return Task{}, fmt.Errorf("not found")
}

type SaveResults struct {
	mu      sync.RWMutex
	Results map[int]float64
}

func (sr *SaveResults) IsExists(id int) bool {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if _, exists := sr.Results[id]; exists {
		return true
	}
	return false
}
func (sr *SaveResults) SetResult(id int, result float64) {
	if Tasks.GetLen() > sr.GetLen() {
		sr.mu.Lock()
		defer sr.mu.Unlock()
		sr.Results[id] = result
		return
	}

}
func (sr *SaveResults) GetResult(id int) float64 {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.Results[id]
}
func (sr *SaveResults) GetLen() int {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return len(sr.Results)
}

var Results = SaveResults{
	mu:      sync.RWMutex{},
	Results: map[int]float64{},
}
var Tasks = SaveTasks{
	mu:    sync.Mutex{},
	Tasks: map[int]Task{},
}

func NewEx(expression string) int {
	id := len(Expressions.Expressions)
	Ex := Expression{Id: id, Expression: strings.ReplaceAll(expression, " ", ""), Status: "processing"}
	Expressions.AddExpression(Ex)
	go func(id int) {
		res, _ := ParseAndEvaluate(Ex)
		Expressions.SetResult(id, res)
		// fmt.Println("Expression ", id, " = ", res)
	}(id)
	return id
}
func ParseAndEvaluate(expression Expression) (float64, error) {
	parser := NewParser(expression.Expression)
	ast := parser.ParseExpression()
	return ast.Evaluate(), nil
}

// Интерфейс для всех узлов AST
type Expr interface {
	Evaluate() float64
}

// Листовой узел для чисел
type Number struct {
	Value float64
}

func (n *Number) Evaluate() float64 {
	return n.Value
}

// Узел для бинарных операций
type BinaryOp struct {
	Left  Expr
	Op    string
	Right Expr
}

func (b *BinaryOp) Evaluate() float64 {
	switch b.Op {
	case "+":
		var res float64
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "+"}
		id := Tasks.GetLen()

		// fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		// fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				// fmt.Printf("res = %f", res)
				// fmt.Println()
				break
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		return res
	case "-":
		var res float64
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "-"}
		id := len(Tasks.Tasks)
		// fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		// fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				// fmt.Printf("res = %f", res)
				break
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		return res
	case "*":
		var res float64
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "*"}
		id := len(Tasks.Tasks)
		// fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		// fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				// fmt.Printf("res = %f", res)
				break
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		return res
	case "/":
		var res float64
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "/"}
		id := len(Tasks.Tasks)
		// fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		// fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				// fmt.Printf("res = %f", res)
				break
			} else {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		return res
	}
	return 0
}

type Parser struct {
	input string
	pos   int
}

func NewParser(input string) *Parser {
	return &Parser{input: input, pos: 0}
}

func (p *Parser) parseNumber() Expr {
	start := p.pos
	for p.pos < len(p.input) && (unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.') {
		p.pos++
	}
	num, _ := strconv.ParseFloat(p.input[start:p.pos], 64)
	return &Number{Value: num}
}

func (p *Parser) ParseExpression() Expr {
	left := p.ParseTerm()

	for p.pos < len(p.input) {
		op := p.input[p.pos]
		if op != '+' && op != '-' {
			break
		}
		p.pos++
		right := p.ParseTerm()
		left = &BinaryOp{Left: left, Op: string(op), Right: right}
	}

	return left
}

func (p *Parser) ParseTerm() Expr {
	left := p.ParseFactor()

	for p.pos < len(p.input) {
		op := p.input[p.pos]
		if op != '*' && op != '/' {
			break
		}
		p.pos++
		right := p.ParseFactor()
		left = &BinaryOp{Left: left, Op: string(op), Right: right}
	}

	return left
}

func (p *Parser) ParseFactor() Expr {
	if p.pos < len(p.input) && p.input[p.pos] == '(' {
		p.pos++ // Пропускаем открывающую скобку
		expr := p.ParseExpression()
		if p.pos < len(p.input) && p.input[p.pos] == ')' {
			p.pos++ // Пропускаем закрывающую скобку
		}
		return expr
	}
	return p.parseNumber()
}

//type User struct {
//	ID             int64
//	Name           string
//	Password       string
//	OriginPassword string
//}
//
//func (u User) ComparePassword(u2 User) error {
//	err := compare(u2.Password, u.OriginPassword)
//	if err != nil {
//		log.Println("auth fail")
//		return err
//	}
//
//	log.Println("auth success")
//	return nil
//}
//
//func createTable(ctx context.Context, db *sql.DB) error {
//	const usersTable = `
//	CREATE TABLE IF NOT EXISTS users(
//		id INTEGER PRIMARY KEY AUTOINCREMENT,
//		name TEXT UNIQUE,
//		password TEXT
//	);`
//
//	if _, err := db.ExecContext(ctx, usersTable); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func insertUser(ctx context.Context, db *sql.DB, user *User) (int64, error) {
//	var q = `
//	INSERT INTO users (name, password) values ($1, $2)
//	`
//	result, err := db.ExecContext(ctx, q, user.Name, user.Password)
//	if err != nil {
//		return 0, err
//	}
//	id, err := result.LastInsertId()
//	if err != nil {
//		return 0, err
//	}
//
//	return id, nil
//}
//
//func selectUser(ctx context.Context, db *sql.DB, name string) (User, error) {
//	var (
//		user User
//		err  error
//	)
//
//	var q = "SELECT id, name, password FROM users WHERE name=$1"
//	err = db.QueryRowContext(ctx, q, name).Scan(&user.ID, &user.Name, &user.Password)
//	return user, err
//}

func generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash := string(hashedBytes[:])
	return hash, nil
}

func compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
func GengerateJWT(user_name string) (string, string) {
	const hmacSampleSecret = "super_secret_signature"
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user_name,
		"nbf":  now.Unix(),
		"exp":  now.Add(15 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user_name,
		"nbf":  now.Unix(),
		"exp":  now.Add(240 * time.Hour).Unix(),
		"iat":  now.Unix(),
	})
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	refresh_tokenString, err := refresh_token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		panic(err)
	}
	fmt.Println(tokenString)
	// return  c.JSON(LoginResponse{AccessToken: refresh_tokenString})
	return tokenString, refresh_tokenString
}

/*
-- Таблица пользователей
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    expression TEXT
);

-- Таблица выражений с статусом вычисления
CREATE TABLE IF NOT EXISTS expression (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expression TEXT NOT NULL,
    answer TEXT,
    userid INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'calculating', 'completed', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userid) REFERENCES user(id) ON DELETE CASCADE
);

-- Индексы для ускорения запросов
CREATE INDEX IF NOT EXISTS idx_expression_userid ON expression(userid);
CREATE INDEX IF NOT EXISTS idx_expression_status ON expression(status);
*/
