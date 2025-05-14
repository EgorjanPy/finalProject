package logic

import (
	"finalProject/internal/storage"
	"finalProject/internal/storage/sqlite"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Expression struct {
	Id         int64
	Expression string
	Result     float64
}

//type SaveExpressions struct {
//	mu          sync.Mutex
//	Expressions []Expression
//}
//
//var Expressions = SaveExpressions{
//	mu:          sync.Mutex{},
//	Expressions: []Expression{},
//}

//func (se *SaveExpressions) GetExpressions() []Expression {
//	se.mu.Lock()
//	defer se.mu.Unlock()
//	return se.Expressions
//}
//
//func (se *SaveExpressions) SetResult(id int, res float64) {
//	se.mu.Lock()
//	se.Expressions[id].Result = res
//	se.Expressions[id].Status = "complited"
//	se.mu.Unlock()
//}
//func (se *SaveExpressions) AddExpression(ex Expression) {
//	se.mu.Lock()
//	se.Expressions = append(se.Expressions, ex)
//	se.mu.Unlock()
//}
//func (se *SaveExpressions) GetExpressionById(id int) (Expression, error) {
//	se.mu.Lock()
//	defer se.mu.Unlock()
//	for _, ex := range se.Expressions {
//		if ex.Id == id {
//			return ex, nil
//		}
//	}
//	return Expression{}, fmt.Errorf("not found")
//}

type Task struct {
	Id        int32
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

func NewEx(expression string) int64 {
	//id := len(Expressions.Expressions)
	id, err := storage.DataBase.AddExpression(&sqlite.Expression{UserID: "1", Expression: expression})
	if err != nil {
		fmt.Printf("Error %v", err)
	}
	Ex := Expression{Id: id, Expression: strings.ReplaceAll(expression, " ", "")}
	go func(id int64) {
		res, _ := ParseAndEvaluate(Ex)
		//Expressions.SetResult(id, res)
		storage.DataBase.SetResult(int64(id), fmt.Sprint(res))
		fmt.Println("Expression ", id, " = ", res)
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

func Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash := string(hashedBytes[:])
	return hash, nil
}
func Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
func ComparePassword(hashedPass, pass string) error {
	err := Compare(hashedPass, pass)
	if err != nil {
		log.Println("auth fail")
		return err
	}

	log.Println("auth success")
	return nil
}

var secretKey = []byte("secret-key")

func CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userID,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	
	return nil
}

//const hmacSampleSecret = "super_secret_signature"
//
//func GengerateJWT(id string) (string, string) {
//	now := time.Now()
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id":  id,
//		"nbf": now.Unix(),
//		"exp": now.Add(24 * time.Hour).Unix(),
//		"iat": now.Unix(),
//	})
//	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"name": id,
//		"nbf":  now.Unix(),
//		"exp":  now.Add(240 * time.Hour).Unix(),
//		"iat":  now.Unix(),
//	})
//	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
//	refresh_tokenString, err := refresh_token.SignedString([]byte(hmacSampleSecret))
//	if err != nil {
//		panic(err)
//	}
//	//fmt.Println(tokenString)
//	// return  c.JSON(LoginResponse{AccessToken: refresh_tokenString})
//	return tokenString, refresh_tokenString
//}

//type User struct {
//	ID             int64
//	Login          string
//	Password       string
//	OriginPassword string
//}

//func ValidateToken(tokenString string) string {
//	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
//		}
//		return []byte(hmacSampleSecret), nil
//	})
//
//	if err != nil {
//		return ""
//	}
//	//return true
//	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
//		return claims["user"]
//	} else {
//		return false
//	}
//}

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
