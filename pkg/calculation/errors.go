package calculation

import "errors"

var (
	ErrEmptyExpression                = errors.New("empty expression")
	ErrInvalidExpression              = errors.New("invalid expression")
	ErrInvalidExpressionInParentheses = errors.New("invalid expression in parentheses")
	ErrDivisionByZero                 = errors.New("division by zero")
	ErrMultiplyError                  = errors.New("multiply error")
	ErrDivisionError                  = errors.New("division error")
)
