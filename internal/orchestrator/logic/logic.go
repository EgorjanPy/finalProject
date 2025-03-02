package logic

import (
	"fmt"
	"strconv"
	"sync"
	"unicode"
)

type Expression struct {
	Id         int
	Expression string
	Status     string
	Result     float64
}

var Expressions = []Expression{}

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

var Tasks = SaveTasks{
	mu:    sync.Mutex{},
	Tasks: map[int]Task{},
}

func (st *SaveTasks) AddTask(id int, task Task) {
	st.mu.Lock()
	st.Tasks[id] = task
	st.mu.Unlock()
}

type SaveResults struct {
	mu      sync.Mutex
	Results map[int]float64
}

func (sr *SaveResults) SetResult(id int, result float64) {
	sr.mu.Lock()

	sr.Results[id] = result
	sr.mu.Unlock()
}
func (sr *SaveResults) GetResult(id int) float64 {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return sr.Results[id]
}

var Results = SaveResults{
	mu:      sync.Mutex{},
	Results: map[int]float64{},
}

func NewEx(expression string) int {
	id := len(Expressions)
	Ex := Expression{Id: id, Expression: expression, Status: "processing"}
	Expressions = append(Expressions, Ex)
	go ParseAndEvaluate(Ex)
	return id
}
func ParseAndEvaluate(expression Expression) (float64, error) {
	parser := NewParser(expression.Expression)
	ast := parser.parseExpression()
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
		newTask := Task{Id: len(Tasks.Tasks), Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "+"}
		id := len(Tasks.Tasks)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.Results[id]; exists {
				res = value
				fmt.Printf("res = %f", res)
				break
			} else {
				// time.Sleep(1 * time.Second)
				continue
			}
		}
		return res
	case "-":
		var res float64
		newTask := Task{Id: len(Tasks.Tasks), Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "-"}
		id := len(Tasks.Tasks)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.Results[id]; exists {
				res = value
				break
			} else {
				continue
			}
		}
		return res
	case "*":
		var res float64
		newTask := Task{Id: len(Tasks.Tasks), Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "*"}
		id := len(Tasks.Tasks)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.Results[id]; exists {
				res = value
				break
			} else {
				continue
			}
		}
		return res
	case "/":
		var res float64
		newTask := Task{Id: len(Tasks.Tasks), Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "/"}
		id := len(Tasks.Tasks)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if value, exists := Results.Results[id]; exists {
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
	return p.parseNumber()
}
