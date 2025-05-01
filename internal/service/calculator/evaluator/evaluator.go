package evaluator

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/parser"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calc_errors"
)

func Evaluate(node *parser.ASTNode) (float64, *calc_errors.EvaluationError) {
	if node == nil {
		return 0, &calc_errors.EvaluationError{
			NodeType: "",
			Reason:   "empty node encountered",
			Context:  "AST structure integrity violation",
		}
	}

	if node.Left == nil && node.Right == nil {
		if node.Token.Type != parser.Number {
			return 0, &calc_errors.EvaluationError{
				NodeType: "character",
				Reason:   "got not a number",
				Context:  fmt.Sprintf("expected number, got %v", node.Token.Value),
			}
		}

		result, err := strconv.ParseFloat(node.Token.Value, 64)
		if err != nil {
			return 0, &calc_errors.EvaluationError{
				NodeType: "character",
				Reason:   "parsing error",
				Context:  fmt.Sprintf("expected number, got %v", node.Token.Value),
			}
		}

		return result, nil
	}

	var leftVal, rightVal float64
	var leftErr, rightErr *calc_errors.EvaluationError

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		leftVal, leftErr = Evaluate(node.Left)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rightVal, rightErr = Evaluate(node.Right)
	}()

	wg.Wait()

	if leftErr != nil {
		return 0, leftErr
	}

	if rightErr != nil {
		return 0, rightErr
	}

	return Calculate(leftVal, rightVal, node.Token.Value)
}
