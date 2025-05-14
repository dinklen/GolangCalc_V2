package evaluator_test

import (
	"testing"
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/evaluator"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calcerr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestCalculate(t *testing.T) {
	cfg := &config.Config{}
	cfg.Microservice.PlusTime = 1 * time.Millisecond
	cfg.Microservice.MinusTime = 1 * time.Millisecond
	cfg.Microservice.DivisionTime = 1 * time.Millisecond
	cfg.Microservice.MultiplyTime = 1 * time.Millisecond
	cfg.Microservice.PowerTime = 1 * time.Millisecond

	logger := zaptest.NewLogger(t)

	testCases := []struct {
		name        string
		left        float64
		right       float64
		operator    string
		expected    float64
		expectedErr error
	}{
		{
			name:     "addition",
			left:     2,
			right:    3,
			operator: "+",
			expected: 5,
		},
		{
			name:     "subtraction",
			left:     5,
			right:    3,
			operator: "-",
			expected: 2,
		},
		{
			name:     "multiplication",
			left:     2,
			right:    3,
			operator: "*",
			expected: 6,
		},
		{
			name:     "division",
			left:     6,
			right:    3,
			operator: "/",
			expected: 2,
		},
		{
			name:        "division by zero",
			left:        5,
			right:       0,
			operator:    "/",
			expectedErr: calcerr.ErrDivisionByZero,
		},
		{
			name:     "power operation",
			left:     2,
			right:    3,
			operator: "^",
			expected: 8,
		},
		{
			name:        "unknown operator",
			left:        2,
			right:       3,
			operator:    "invalid",
			expectedErr: calcerr.ErrUnknownOperator,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			result, err := evaluator.Calculate(tc.left, tc.right, tc.operator, cfg, logger)
			duration := time.Since(start)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)

			minDuration := getMinDuration(tc.operator, cfg)
			assert.True(t, duration >= minDuration, "operation took %s, expected at least %s", duration, minDuration)
		})
	}
}

func getMinDuration(operator string, cfg *config.Config) time.Duration {
	switch operator {
	case "+":
		return cfg.Microservice.PlusTime
	case "-":
		return cfg.Microservice.MinusTime
	case "*":
		return cfg.Microservice.MultiplyTime
	case "/":
		return cfg.Microservice.DivisionTime
	case "^":
		return cfg.Microservice.PowerTime
	default:
		return 0
	}
}
