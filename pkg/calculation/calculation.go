package calculation

import (
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	if expression == "" {
		return 0, ErrEmptyExpression
	}

	var nums []float64
	var operations []rune

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	applyOp := func() error {
		if len(nums) < 2 || len(operations) == 0 {
			return ErrInvalidExpression
		}
		b := nums[len(nums)-1]
		a := nums[len(nums)-2]
		op := operations[len(operations)-1]
		nums = nums[:len(nums)-2]
		operations = operations[:len(operations)-1]

		var result float64
		switch op {
		case '+':
			result = a + b
		case '-':
			result = a - b
		case '*':
			result = a * b
		case '/':
			if b == 0 {
				return ErrDivisionByZero
			}
			result = a / b
		default:
			return ErrInvalidExpression
		}
		nums = append(nums, result)
		return nil
	}

	i := 0
	for i < len(expression) {
		if expression[i] >= '0' && expression[i] <= '9' {
			start := i
			for i < len(expression) && (expression[i] >= '0' && expression[i] <= '9' || expression[i] == '.') {
				i++
			}
			num, err := strconv.ParseFloat(expression[start:i], 64)
			if err != nil {
				return 0, err
			}
			nums = append(nums, num)
			continue
		}

		if expression[i] == '(' {
			operations = append(operations, '(')
		} else if expression[i] == ')' {
			for len(operations) > 0 && operations[len(operations)-1] != '(' {
				if err := applyOp(); err != nil {
					return 0, err
				}
			}
			if len(operations) == 0 {
				return 0, ErrInvalidExpressionInParentheses
			}
			operations = operations[:len(operations)-1]
		} else {
			for len(operations) > 0 && precedence[rune(expression[i])] <= precedence[operations[len(operations)-1]] {
				if err := applyOp(); err != nil {
					return 0, err
				}
			}
			operations = append(operations, rune(expression[i]))
		}
		i++
	}
	for len(operations) > 0 {
		if err := applyOp(); err != nil {
			return 0, err
		}
	}
	if len(nums) != 1 {
		return 0, ErrInvalidExpression
	}

	return nums[0], nil
}
