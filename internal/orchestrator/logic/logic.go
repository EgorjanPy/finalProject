package logic

import (
	"fmt"
	"strconv"
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

func (se *SaveExpressions) SetResult(id int, res float64) {
	se.mu.Lock()
	se.Expressions[id].Result = res
	se.mu.Unlock()
}
func (se *SaveExpressions) AddExpression(ex Expression) {
	se.mu.Lock()
	se.Expressions = append(se.Expressions, ex)
	se.mu.Unlock()
}
func (se *SaveExpressions) GetExpressionById(id int) Expression {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.Expressions[id]
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
func (st *SaveTasks) GetTaskById(id int) Task {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.Tasks[id]
}

type SaveResults struct {
	mu      sync.RWMutex
	Results map[int]float64
}

func (sr *SaveResults) IsExists(id int) bool {
	if _, exists := Results.Results[id]; exists {
		return true
	}
	return false
}
func (sr *SaveResults) SetResult(id int, result float64) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.Results[id] = result
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
	Ex := Expression{Id: id, Expression: expression, Status: "processing"}
	Expressions.AddExpression(Ex)
	go func(id int) {
		res, _ := ParseAndEvaluate(Ex)
		Expressions.SetResult(id, res)
		fmt.Println("Expression ", id, " = ", res)
	}(id)
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
		newTask := Task{Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "+"}
		id := Tasks.GetLen()

		fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				fmt.Printf("res = %f", res)
				fmt.Println()
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
		fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				fmt.Printf("res = %f", res)
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
		fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				fmt.Printf("res = %f", res)
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
		fmt.Println("len, id = ", id)
		Tasks.AddTask(id, newTask)
		fmt.Println(Tasks.Tasks)
		for {
			if Results.IsExists(id) {
				res = Results.GetResult(id)
				fmt.Printf("res = %f", res)
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

// package logic

// import (
// 	"fmt"
// 	"strconv"
// 	"sync"
// 	"time"
// 	"unicode"
// )

// // Expression представляет математическое выражение
// type Expression struct {
// 	Id         int
// 	Expression string
// 	Status     string
// 	Result     float64
// }
// type SaveExpressions struct {
// 	mu          sync.Mutex
// 	Expressions []Expression
// }

// func (se *SaveExpressions) AddExpression(ex Expression) {
// 	se.mu.Lock()
// 	se.Expressions = append(se.Expressions, ex)
// 	se.mu.Unlock()
// }
// func (se *SaveExpressions) GetLen() int {
// 	se.mu.Lock()
// 	defer se.mu.Unlock()
// 	return len(se.Expressions)
// }
// func (se *SaveExpressions) GetExpressionById(id int) Expression {
// 	se.mu.Lock()
// 	defer se.mu.Unlock()
// 	return se.Expressions[id]
// }

// var Expressions = SaveExpressions{
// 	mu:          sync.Mutex{},
// 	Expressions: []Expression{},
// }

// // Task представляет задачу для выполнения операции
// type Task struct {
// 	Id        int
// 	Arg1      float64
// 	Arg2      float64
// 	Operation string
// }

// // SaveTasks хранит задачи
// type SaveTasks struct {
// 	mu    sync.Mutex
// 	Tasks map[int]Task
// }

// var Tasks = SaveTasks{
// 	mu:    sync.Mutex{},
// 	Tasks: map[int]Task{},
// }

// // AddTask добавляет задачу в хранилище
// func (st *SaveTasks) AddTask(id int, task Task) {
// 	st.mu.Lock()
// 	st.Tasks[id] = task
// 	st.mu.Unlock()
// }
// func (st *SaveTasks) GetLen() int {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	return len(st.Tasks)
// }
// func (st *SaveTasks) GetTaskById(id int) Task {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	return st.Tasks[id]
// }

// // SaveResults хранит результаты вычислений
// type SaveResults struct {
// 	mu      sync.Mutex
// 	Results map[int]float64
// }

// // SetResult устанавливает результат для задачи
// func (sr *SaveResults) SetResult(id int, result float64) {
// 	sr.mu.Lock()
// 	sr.Results[id] = result
// 	sr.mu.Unlock()
// }
// func (sr *SaveResults) GetLen() int {
// 	sr.mu.Lock()
// 	defer sr.mu.Unlock()
// 	return len(sr.Results)
// }

// // GetResultById получает результат по ID
// func (sr *SaveResults) IsExists(id int) bool {
// 	if _, exists := Results.Results[id]; exists {
// 		fmt.Println("Да существет")
// 		return true
// 	}
// 	fmt.Println("Нет не существет")
// 	return false
// }
// func (sr *SaveResults) GetResultById(id int) float64 {
// 	sr.mu.Lock()
// 	defer sr.mu.Unlock()
// 	return Results.Results[id]
// }

// var Results = SaveResults{
// 	mu:      sync.Mutex{},
// 	Results: map[int]float64{},
// }

// // NewEx создает новое выражение и запускает его обработку
// func NewEx(expression string) int {
// 	id := Expressions.GetLen()
// 	ex := Expression{Id: id, Expression: expression, Status: "processing"}
// 	Expressions.AddExpression(ex)
// 	ParseAndEvaluate(ex)
// 	return id
// }
// func ParseAndEvaluate(expression Expression) (float64, error) {
// 	parser := NewParser(expression.Expression)
// 	ast := parser.parseExpression()
// 	return ast.Evaluate(), nil
// }

// // Интерфейс для всех узлов AST
// type Expr interface {
// 	Evaluate() float64
// }

// // Листовой узел для чисел
// type Number struct {
// 	Value float64
// }

// func (n *Number) Evaluate() float64 {
// 	return n.Value
// }

// // Узел для бинарных операций
// type BinaryOp struct {
// 	Left  Expr
// 	Op    string
// 	Right Expr
// }

// func (b *BinaryOp) Evaluate() float64 {
// 	switch b.Op {
// 	case "+":
// 		var res float64
// 		id := Tasks.GetLen()
// 		fmt.Printf("id = %d", id)
// 		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "+"}
// 		Tasks.AddTask(id, newTask)
// 		fmt.Println(Tasks.Tasks)
// 		for {
// 			if Results.IsExists(id) {
// 				res = Results.GetResultById(id)
// 				fmt.Println("Tyt")
// 				fmt.Println(res)
// 				break
// 			} else {
// 				time.Sleep(time.Second * 1)
// 				continue
// 			}
// 		}
// 		return res
// 	case "-":
// 		var res float64
// 		id := Tasks.GetLen()
// 		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "-"}
// 		Tasks.AddTask(id, newTask)
// 		fmt.Println(Tasks.Tasks)
// 		for {
// 			if Results.IsExists(id) {
// 				res = Results.GetResultById(id)
// 				fmt.Println("Tyt")
// 				break
// 			} else {
// 				time.Sleep(time.Second * 1)
// 				continue
// 			}
// 		}
// 		return res
// 	case "*":
// 		var res float64
// 		id := Tasks.GetLen()
// 		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "*"}
// 		Tasks.AddTask(id, newTask)
// 		fmt.Println(Tasks.Tasks)
// 		for {
// 			if Results.IsExists(id) {
// 				res = Results.GetResultById(id)
// 				fmt.Println("Tyt")
// 				break
// 			} else {
// 				time.Sleep(time.Second * 1)
// 				continue
// 			}
// 		}
// 		return res
// 	case "/":
// 		var res float64
// 		id := Tasks.GetLen()
// 		newTask := Task{Id: id, Arg1: b.Left.Evaluate(), Arg2: b.Right.Evaluate(), Operation: "/"}
// 		Tasks.AddTask(id, newTask)
// 		fmt.Println(Tasks.Tasks)
// 		for {
// 			if Results.IsExists(id) {
// 				res = Results.GetResultById(id)
// 				fmt.Println("Tyt")
// 				break
// 			} else {
// 				time.Sleep(time.Second * 1)
// 				continue
// 			}
// 		}
// 		return res
// 	}
// 	return 0
// }

// type Parser struct {
// 	input string
// 	pos   int
// }

// func NewParser(input string) *Parser {
// 	return &Parser{input: input, pos: 0}
// }

// func (p *Parser) parseNumber() Expr {
// 	start := p.pos
// 	for p.pos < len(p.input) && (unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.') {
// 		p.pos++
// 	}
// 	num, _ := strconv.ParseFloat(p.input[start:p.pos], 64)
// 	return &Number{Value: num}
// }

// func (p *Parser) parseExpression() Expr {
// 	left := p.parseTerm()

// 	for p.pos < len(p.input) {
// 		op := p.input[p.pos]
// 		if op != '+' && op != '-' {
// 			break
// 		}
// 		p.pos++
// 		right := p.parseTerm()
// 		left = &BinaryOp{Left: left, Op: string(op), Right: right}
// 	}

// 	return left
// }

// func (p *Parser) parseTerm() Expr {
// 	left := p.parseFactor()
// 	for p.pos < len(p.input) {
// 		op := p.input[p.pos]
// 		if op != '*' && op != '/' {
// 			break
// 		}
// 		p.pos++
// 		right := p.parseFactor()
// 		left = &BinaryOp{Left: left, Op: string(op), Right: right}
// 	}

// 	return left
// }

// func (p *Parser) parseFactor() Expr {
// 	return p.parseNumber()
// }
