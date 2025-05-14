package parser

import (
	"testing"

	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calcerr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBuildAST_SimpleExpression(t *testing.T) {
	tokens := []Token{
		{Type: Number, Value: "2"},
		{Type: Operator, Value: "+"},
		{Type: Number, Value: "3"},
	}

	ast, err := BuildAST(tokens, zap.NewNop())

	assert.NoError(t, err)
	assert.Equal(t, "+", ast.Token.Value)
	assert.Equal(t, "2", ast.Left.Token.Value)
	assert.Equal(t, "3", ast.Right.Token.Value)
}

func TestBuildAST_OperatorPrecedence(t *testing.T) {
	tokens := []Token{
		{Type: Number, Value: "2"},
		{Type: Operator, Value: "+"},
		{Type: Number, Value: "3"},
		{Type: Operator, Value: "*"},
		{Type: Number, Value: "4"},
	}

	ast, err := BuildAST(tokens, zap.NewNop())

	assert.NoError(t, err)
	assert.Equal(t, "+", ast.Token.Value)
	assert.Equal(t, "*", ast.Right.Token.Value)
	assert.Equal(t, "3", ast.Right.Left.Token.Value)
	assert.Equal(t, "4", ast.Right.Right.Token.Value)
}

func TestBuildAST_Parentheses(t *testing.T) {
	tokens := []Token{
		{Type: LeftParen, Value: "("},
		{Type: Number, Value: "2"},
		{Type: Operator, Value: "+"},
		{Type: Number, Value: "3"},
		{Type: RightParen, Value: ")"},
		{Type: Operator, Value: "*"},
		{Type: Number, Value: "4"},
	}

	ast, err := BuildAST(tokens, zap.NewNop())

	assert.NoError(t, err)
	assert.Equal(t, "*", ast.Token.Value)
	assert.Equal(t, "+", ast.Left.Token.Value)
	assert.Equal(t, "4", ast.Right.Token.Value)
}

func TestBuildAST_MismatchedParentheses(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
	}{
		{
			name: "unclosed parenthesis",
			tokens: []Token{
				{Type: LeftParen, Value: "("},
				{Type: Number, Value: "2"},
				{Type: Operator, Value: "+"},
				{Type: Number, Value: "3"},
			},
		},
		{
			name: "unopened parenthesis",
			tokens: []Token{
				{Type: Number, Value: "2"},
				{Type: RightParen, Value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildAST(tt.tokens, zap.NewNop())
			assert.ErrorIs(t, err, calcerr.ErrMismatchedParentheses)
		})
	}
}

func TestBuildAST_NotEnoughOperands(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
	}{
		{
			name: "single operator",
			tokens: []Token{
				{Type: Operator, Value: "+"},
			},
		},
		{
			name: "operator without operands",
			tokens: []Token{
				{Type: Number, Value: "2"},
				{Type: Operator, Value: "+"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildAST(tt.tokens, zap.NewNop())
			assert.ErrorIs(t, err, calcerr.ErrNotEnoughOperands)
		})
	}
}

func TestBuildAST_InvalidExpression(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
	}{
		{
			name: "too many numbers",
			tokens: []Token{
				{Type: Number, Value: "2"},
				{Type: Number, Value: "3"},
			},
		},
		{
			name: "multiple roots",
			tokens: []Token{
				{Type: Number, Value: "2"},
				{Type: Operator, Value: "+"},
				{Type: Number, Value: "3"},
				{Type: Operator, Value: "+"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildAST(tt.tokens, zap.NewNop())
			assert.ErrorIs(t, err, calcerr.ErrInvalidExpression)
		})
	}
}
