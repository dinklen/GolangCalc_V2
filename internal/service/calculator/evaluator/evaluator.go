package evaluator

import (
	"fmt"
	"strconv"
	"sync"
	"context"

	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/parser"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calc_errors"
	"github.com/dinklen/GolangCalc_V2/internal/task_manager"
)

func Evaluate(node *parser.ASTNode,tm *task_manager.TaskManager, parentID string) (float64, *calc_errors.EvaluationError) {
	if node == nil {
		return 0.0, &calc_errors.EvaluationError{
			NodeType: "",
			Reason:   "empty node encountered",
			Context:  "AST structure integrity violation",
		}
	}

	if node.Left == nil && node.Right == nil {
		if node.Token.Type != parser.Number {
			return 0.0, &calc_errors.EvaluationError{
				NodeType: "character",
				Reason:   "got not a number",
				Context:  fmt.Sprintf("expected number, got %v", node.Token.Value),
			}
		}

		result, err := strconv.ParseFloat(node.Token.Value, 64)
		if err != nil {
			return 0.0, &calc_errors.EvaluationError{
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
		return 0.0, leftErr
	}

	if rightErr != nil {
		return 0.0, rightErr
	}

	stream, err := client.Calculate(context.Background())
	if err != nil {
		fmt.Errorf("Failed to create stream: %v", err)
	}

	result, err := tm.Evaluate(leftVal, rightVal, node.Token.Value, parentID)
	if err != nil {
		return 0.0, &calc_errors.EvaluationError{
			NodeType: "none",
			Reason: "microservice",
			Context: fmt.Sprintf("%v", err),
		}
	}
	return result, nil
}
