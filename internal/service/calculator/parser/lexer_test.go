package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTokenize_BasicOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple addition",
			input: "3+4",
			expected: []Token{
				{Type: Number, Value: "3"},
				{Type: Operator, Value: "+"},
				{Type: Number, Value: "4"},
			},
		},
		{
			name:  "decimal numbers",
			input: "12.34 - 5.6",
			expected: []Token{
				{Type: Number, Value: "12.34"},
				{Type: Operator, Value: "-"},
				{Type: Number, Value: "5.6"},
			},
		},
		{
			name:  "multiple operators",
			input: "2*3^4",
			expected: []Token{
				{Type: Number, Value: "2"},
				{Type: Operator, Value: "*"},
				{Type: Number, Value: "3"},
				{Type: Operator, Value: "^"},
				{Type: Number, Value: "4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Tokenize(tt.input, zap.NewNop())
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenize_Parentheses(t *testing.T) {
	input := "(2+3)*4"
	expected := []Token{
		{Type: LeftParen, Value: "("},
		{Type: Number, Value: "2"},
		{Type: Operator, Value: "+"},
		{Type: Number, Value: "3"},
		{Type: RightParen, Value: ")"},
		{Type: Operator, Value: "*"},
		{Type: Number, Value: "4"},
	}

	result, err := Tokenize(input, zap.NewNop())
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTokenize_UnaryMinus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "unary at start",
			input: "-5",
			expected: []Token{
				{Type: Number, Value: "0"},
				{Type: Operator, Value: "-"},
				{Type: Number, Value: "5"},
			},
		},
		{
			name:  "unary after operator",
			input: "2*(-3)",
			expected: []Token{
				{Type: Number, Value: "2"},
				{Type: Operator, Value: "*"},
				{Type: LeftParen, Value: "("},
				{Type: Number, Value: "0"},
				{Type: Operator, Value: "-"},
				{Type: Number, Value: "3"},
				{Type: RightParen, Value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Tokenize(tt.input, zap.NewNop())
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenize_WhitespaceHandling(t *testing.T) {
	input := "  2  +  \t 3.5 "
	expected := []Token{
		{Type: Number, Value: "2"},
		{Type: Operator, Value: "+"},
		{Type: Number, Value: "3.5"},
	}

	result, err := Tokenize(input, zap.NewNop())
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
