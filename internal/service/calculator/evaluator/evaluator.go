package evaluator

import (
	"strconv"
	"sync"

	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/parser"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calcerr"
	"github.com/dinklen/GolangCalc_V2/internal/tm"

	"go.uber.org/zap"
)

func Evaluate(node *parser.ASTNode, taskManager *tm.TaskManager, parentID string, logger *zap.Logger) (float64, error) {
	if node == nil {
		return 0.0, calcerr.ErrEmtpyNode
	}
	logger.Info("subtree evaluation started")

	if node.Left == nil && node.Right == nil {
		if node.Token.Type != parser.Number {
			return 0.0, calcerr.ErrUnexpectedCharacter
		}
		logger.Info("number node found")

		result, err := strconv.ParseFloat(node.Token.Value, 64)
		if err != nil {
			return 0.0, calcerr.ErrParsingFailed
		}
		logger.Info("number node parsing success")

		return result, nil
	}

	var leftVal, rightVal float64
	var leftErr, rightErr error

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		leftVal, leftErr = Evaluate(node.Left, taskManager, parentID, logger)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rightVal, rightErr = Evaluate(node.Right, taskManager, parentID, logger)
	}()

	wg.Wait()

	if leftErr != nil {
		return 0.0, leftErr
	}

	if rightErr != nil {
		return 0.0, rightErr
	}
	logger.Info("subtrees evaluation finished successful")

	result, err := taskManager.Evaluate(leftVal, rightVal, node.Token.Value, parentID, logger)
	if err != nil {
		return 0.0, calcerr.ErrBadAnswerFromMicroservice
	}

	logger.Info("subtree evaluation finished")
	return result, nil
}
