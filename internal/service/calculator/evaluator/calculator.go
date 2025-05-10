package evaluator

import (
	"fmt"

	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calc_errors"
)

func Calculate(leftVal, rightVal float64, operator string) (float64, *calc_errors.EvaluationError) {
	switch operator {
	case "+":
		return leftVal + rightVal, nil

	case "-":
		return leftVal - rightVal, nil

	case "*":
		return leftVal * rightVal, nil

	case "/":
		if rightVal == 0 {
			return 0, &calc_errors.EvaluationError{
				NodeType: "operator",
				Reason:   "division by zero",
				Context:  "attempt to divide by zero",
			}
		}

		return leftVal / rightVal, nil

	case "^":
		result := 1.0
		for i := 0; i < int(rightVal); i++ {
			result *= leftVal
		}

		return result, nil

	default:
		return 0, &calc_errors.EvaluationError{
			NodeType: "operator",
			Reason:   "unknown operator",
			Context:  fmt.Sprintf("failed to recognize operator \"%s\"", operator),
		}
	}
}

// Задержки из конфига
