package calculation

import "errors"

var (
	ErrEmptyExpression   = errors.New("empty expression")
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero    = errors.New("division by zero")
	// ErrUnknowError       = errors.New("unknow error")
	// ErrMultiplyError     = errors.New("multiply error")
	// ErrDivisionError     = errors.New("division error")
)
