package logic_test

// import (
// 	"finalProject/internal/orchestrator/logic"
// 	"testing"
// )

// // Expressions
// func TestAddExpression(t *testing.T) {
// 	exp := logic.Expression{Id: 1, Expression: "2+2", Status: "processing"}
// 	logic.Expressions.AddExpression(exp)

// 	if len(logic.Expressions.GetExpressions()) != 1 {
// 		t.Errorf("Expected 1 expression, got %d", len(logic.Expressions.GetExpressions()))
// 	}
// }

// // Tasks
// func TestSetResult(t *testing.T) {
// 	exp := logic.Expression{Id: 1, Expression: "2+2", Status: "processing"}
// 	logic.Expressions.AddExpression(exp)
// 	logic.Expressions.SetResult(1, 4.0)

// 	result, err := logic.Expressions.GetExpressionById(1)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if result.Result != 4.0 || result.Status != "complited" {
// 		t.Errorf("Expected result 4.0 and status 'complited', got %f and %s", result.Result, result.Status)
// 	}
// }

// func TestGetExpressionById(t *testing.T) {
// 	exp := logic.Expression{Id: 1, Expression: "2+2", Status: "processing"}
// 	logic.Expressions.AddExpression(exp)

// 	result, err := logic.Expressions.GetExpressionById(1)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if result.Id != 1 {
// 		t.Errorf("Expected ID 1, got %d", result.Id)
// 	}
// }
// func TestAddTask(t *testing.T) {
// 	task := logic.Task{Id: 1, Arg1: 2, Arg2: 2, Operation: "+"}
// 	logic.Tasks.AddTask(1, task)

// 	if len(logic.Tasks.Tasks) != 1 {
// 		t.Errorf("Expected 1 task, got %d", len(logic.Tasks.Tasks))
// 	}
// }

// func TestGetTaskById(t *testing.T) {
// 	task := logic.Task{Id: 1, Arg1: 2, Arg2: 2, Operation: "+"}
// 	logic.Tasks.AddTask(1, task)

// 	result, err := logic.Tasks.GetTaskById(1)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if result.Id != 1 {
// 		t.Errorf("Expected ID 1, got %d", result.Id)
// 	}
// }

// // Result

// func TestGetResult(t *testing.T) {
// 	logic.Results.SetResult(1, 4.0)

// 	result := logic.Results.GetResult(1)
// 	if result != 4.0 {
// 		t.Errorf("Expected result 4.0, got %f", result)
// 	}
// }

// // Parser
