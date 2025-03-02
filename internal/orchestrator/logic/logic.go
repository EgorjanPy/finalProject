package logic

import (
	"fmt"
	"strconv"
	"sync"
	"unicode"
)

// Expression представляет математическое выражение
type Expression struct {
	Id         int
	Expression string
	Status     string
	Result     float64
}

var Expressions = []Expression{}

// Task представляет задачу для выполнения операции
type Task struct {
	Id        int
	Arg1      float64
	Arg2      float64
	Operation string
}

// SaveTasks хранит задачи
type SaveTasks struct {
	mu    sync.Mutex
	Tasks map[int]Task
}

var Tasks = SaveTasks{
	mu:    sync.Mutex{},
	Tasks: map[int]Task{},
}

// AddTask добавляет задачу в хранилище
func (st *SaveTasks) AddTask(id int, task Task) {
	st.mu.Lock()
	st.Tasks[id] = task
	st.mu.Unlock()
}
func (st *SaveTasks) GetLen() int {
	st.mu.Lock()
	defer st.mu.Unlock()
	return len(st.Tasks)
}

// SaveResults хранит результаты вычислений
type SaveResults struct {
	mu      sync.Mutex
	Results map[int]float64
}

// SetResult устанавливает результат для задачи
func (sr *SaveResults) SetResult(id int, result float64) {
	sr.mu.Lock()
	sr.Results[id] = result
	sr.mu.Unlock()
}
func (sr *SaveResults) GetLen() int {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return len(sr.Results)
}

// GetResultById получает результат по ID
func (sr *SaveResults) GetResultById(id int) (float64, bool) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if value, exists := Results.Results[id]; exists {
		return value, true
	} else {
		return 0, false
	}
}

var Results = SaveResults{
	mu:      sync.Mutex{},
	Results: map[int]float64{},
}

// NewEx создает новое выражение и запускает его обработку
func NewEx(expression string) int {
	id := len(Expressions)
	Ex := Expression{Id: id, Expression: expression, Status: "processing"}
	Expressions = append(Expressions, Ex)
	go ParseAndEvaluate(Ex)
	return id
}

// ParseAndEvaluate парсит и вычисляет выражение
func ParseAndEvaluate(expression Expression) (float64, error) {
	parser := NewParser(expression.Expression)
	ast := parser.parseExpression()
	result := ast.Evaluate()
	Results.SetResult(expression.Id, result)
	// expression.Result = result
	// expression.Status = "Complited"
	return result, nil
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
		id := Tasks.GetLen()
		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "+"}
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.GetResultById(id); exists {
				res = value
				break
			} else {
				continue
			}
		}
		return res
	case "-":
		var res float64
		id := Tasks.GetLen()
		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "-"}
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.GetResultById(id); exists {
				res = value
				break
			} else {
				continue
			}
		}
		return res
	case "*":
		var res float64
		id := Tasks.GetLen()
		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "*"}
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.GetResultById(id); exists {
				res = value
				break
			} else {
				continue
			}
		}
		return res
	case "/":
		var res float64
		id := Tasks.GetLen()
		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "/"}
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.GetResultById(id); exists {
				res = value
				break
			} else {
				continue
			}
		}
		return res
	}
	return 0
}

// Parser для математических выражений
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

func (p *Parser) parseExpression() Expr {
	left := p.parseTerm()

	for p.pos < len(p.input) {
		op := p.input[p.pos]
		if op != '+' && op != '-' {
			break
		}
		p.pos++
		right := p.parseTerm()
		left = &BinaryOp{Left: left, Op: string(op), Right: right}
	}

	return left
}

func (p *Parser) parseTerm() Expr {
	left := p.parseFactor()

	for p.pos < len(p.input) {
		op := p.input[p.pos]
		if op != '*' && op != '/' {
			break
		}
		p.pos++
		right := p.parseFactor()
		left = &BinaryOp{Left: left, Op: string(op), Right: right}
	}

	return left
}

func (p *Parser) parseFactor() Expr {
	if p.pos < len(p.input) && p.input[p.pos] == '(' {
		p.pos++ // Пропускаем открывающую скобку
		expr := p.parseExpression()
		if p.pos < len(p.input) && p.input[p.pos] == ')' {
			p.pos++ // Пропускаем закрывающую скобку
		}
		return expr
	}
	return p.parseNumber()
}
