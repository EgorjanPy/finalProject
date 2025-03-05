package logic

import (
	"testing"
)

// Expressions
func TestAddExpression(t *testing.T) {
	exp := Expression{Id: 1, Expression: "2+2", Status: "processing"}
	Expressions.AddExpression(exp)

	if len(Expressions.GetExpressions()) != 1 {
		t.Errorf("Expected 1 expression, got %d", len(Expressions.GetExpressions()))
	}
}

// Tasks
func TestSetResult(t *testing.T) {
	exp := Expression{Id: 1, Expression: "2+2", Status: "processing"}
	Expressions.AddExpression(exp)
	Expressions.SetResult(1, 4.0)

	result, err := Expressions.GetExpressionById(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Result != 4.0 || result.Status != "complited" {
		t.Errorf("Expected result 4.0 and status 'complited', got %f and %s", result.Result, result.Status)
	}
}

func TestGetExpressionById(t *testing.T) {
	exp := Expression{Id: 1, Expression: "2+2", Status: "processing"}
	Expressions.AddExpression(exp)

	result, err := Expressions.GetExpressionById(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Id != 1 {
		t.Errorf("Expected ID 1, got %d", result.Id)
	}
}
func TestAddTask(t *testing.T) {
	task := Task{Id: 1, Arg1: 2, Arg2: 2, Operation: "+"}
	Tasks.AddTask(1, task)

	if len(Tasks.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(Tasks.Tasks))
	}
}

func TestGetTaskById(t *testing.T) {
	task := Task{Id: 1, Arg1: 2, Arg2: 2, Operation: "+"}
	Tasks.AddTask(1, task)

	result, err := Tasks.GetTaskById(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Id != 1 {
		t.Errorf("Expected ID 1, got %d", result.Id)
	}
}

// Result

func TestGetResult(t *testing.T) {
	Results.SetResult(1, 4.0)

	result := Results.GetResult(1)
	if result != 4.0 {
		t.Errorf("Expected result 4.0, got %f", result)
	}
}

// Parser
func TestParseNumber(t *testing.T) {
	p := NewParser("123.45")
	num := p.parseNumber().(*Number)

	if num.Value != 123.45 {
		t.Errorf("Expected 123.45, got %f", num.Value)
	}
}

func TestParseExpression(t *testing.T) {
	p := NewParser("2+3*4")
	expr := p.ParseExpression()

	result := expr.Evaluate()
	if result != 14.0 {
		t.Errorf("Expected 14.0, got %f", result)
	}
}

func TestParseTerm(t *testing.T) {
	p := NewParser("2*3+4")
	term := p.ParseTerm()

	result := term.Evaluate()
	if result != 6.0 {
		t.Errorf("Expected 6.0, got %f", result)
	}
}

func TestParseFactor(t *testing.T) {
	p := NewParser("(2+3)*4")
	factor := p.ParseFactor()

	result := factor.Evaluate()
	if result != 20.0 {
		t.Errorf("Expected 20.0, got %f", result)
	}
}
