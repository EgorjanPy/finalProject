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
type Task struct {
	Id        int
	Arg1      float64
	Arg2      float64
	Operation string
	res       float64
}

var Expressions = []Expression{}

// var Tasks = map[int]Task{}
type SaveTasks struct {
	mu    sync.Mutex
	Tasks map[int]Task
}
type SaveResults struct {
	mu      sync.RWMutex
	Results map[int]float64
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
	id := len(Expressions)
	Ex := Expression{Id: id, Expression: expression, Status: "processing"}
	Expressions = append(Expressions, Ex)
	ParseAndEvaluate(Ex)
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
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "+"}
		Tasks.mu.Lock()
		id := len(Tasks.Tasks)
		Tasks.Tasks[id] = newTask
		fmt.Println(Tasks.Tasks)
		Tasks.mu.Unlock()
		Results.mu.RLock()
		res := Results.Results[id]
		Results.mu.RUnlock()
		return res
	case "-":
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "-"}
		Tasks.mu.Lock()
		id := len(Tasks.Tasks)
		Tasks.Tasks[id] = newTask
		Tasks.mu.Unlock()
		Results.mu.RLock()
		res := Results.Results[id]
		Results.mu.RUnlock()
		return res
	case "*":
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "*"}
		Tasks.mu.Lock()
		id := len(Tasks.Tasks)
		Tasks.Tasks[id] = newTask
		Tasks.mu.Unlock()
		Results.mu.RLock()
		res := Results.Results[id]
		Results.mu.RUnlock()
		return res
	case "/":
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "/"}
		Tasks.mu.Lock()
		id := len(Tasks.Tasks)
		Tasks.Tasks[id] = newTask
		Tasks.mu.Unlock()
		Results.mu.RLock()
		res := Results.Results[id]
		Results.mu.RUnlock()
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
