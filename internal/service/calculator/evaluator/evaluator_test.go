package evaluator

import (
	"testing"

	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/parser"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calcerr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockTaskManager struct {
	result   float64
	err      error
	callArgs []interface{}
}

func (m *MockTaskManager) Evaluate(left, right float64, op, parentID string, logger *zap.Logger) (float64, error) {
	m.callArgs = []interface{}{left, right, op, parentID}
	return m.result, m.err
}

func TestEvaluate_NilNode(t *testing.T) {
	result, err := Evaluate(nil, nil, "", zap.NewNop())
	assert.ErrorIs(t, err, calcerr.ErrEmtpyNode)
	assert.Equal(t, 0.0, result)
}

func TestEvaluate_NumberNodeSuccess(t *testing.T) {
	node := &parser.ASTNode{
		Token: parser.Token{Type: parser.Number, Value: "42.5"},
	}
	result, err := Evaluate(node, nil, "", zap.NewNop())
	assert.NoError(t, err)
	assert.Equal(t, 42.5, result)
}

func TestEvaluate_NumberNodeInvalid(t *testing.T) {
	node := &parser.ASTNode{
		Token: parser.Token{Type: parser.Number, Value: "invalid"},
	}
	result, err := Evaluate(node, nil, "", zap.NewNop())
	assert.ErrorIs(t, err, calcerr.ErrParsingFailed)
	assert.Equal(t, 0.0, result)
}
