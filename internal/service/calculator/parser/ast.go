package parser

import (
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calcerr"

	"go.uber.org/zap"
)

type ASTNode struct {
	Token   Token
	Left    *ASTNode
	Right   *ASTNode
	IsUnary bool
}

func BuildAST(tokens []Token, logger *zap.Logger) (*ASTNode, error) {
	logger.Info("build AST started")
	var output []*ASTNode
	var operators []Token

	for _, token := range tokens {
		switch token.Type {
		case Number:
			output = append(output, &ASTNode{Token: token})

		case Operator:
			for len(operators) > 0 {
				top := operators[len(operators)-1]
				if top.Type != Operator || precedence(token.Value) > precedence(top.Value) {
					break
				}

				operators = operators[:len(operators)-1]
				if len(output) < 2 {
					return nil, calcerr.ErrNotEnoughOperands
				}

				right := output[len(output)-1]
				left := output[len(output)-2]
				output = output[:len(output)-2]

				output = append(output, &ASTNode{
					Token: top,
					Left:  left,
					Right: right,
				})
			}

			operators = append(operators, token)

		case LeftParen:
			operators = append(operators, token)

		case RightParen:
			for len(operators) > 0 && operators[len(operators)-1].Type != LeftParen {
				top := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				if len(output) < 2 {
					return nil, calcerr.ErrNotEnoughOperands
				}

				right := output[len(output)-1]
				left := output[len(output)-2]
				output = output[:len(output)-2]

				output = append(output, &ASTNode{
					Token: top,
					Left:  left,
					Right: right,
				})
			}

			if len(operators) == 0 || operators[len(operators)-1].Type != LeftParen {
				return nil, calcerr.ErrMismatchedParentheses
			}

			operators = operators[:len(operators)-1]
		}
	}

	for len(operators) > 0 {
		top := operators[len(operators)-1]
		if top.Type == LeftParen {
			return nil, calcerr.ErrMismatchedParentheses
		}

		operators = operators[:len(operators)-1]
		if len(output) < 2 {
			return nil, calcerr.ErrNotEnoughOperands
		}

		right := output[len(output)-1]
		left := output[len(output)-2]
		output = output[:len(output)-2]

		output = append(output, &ASTNode{
			Token: top,
			Left:  left,
			Right: right,
		})
	}

	if len(output) != 1 {
		return nil, calcerr.ErrInvalidExpression
	}

	logger.Info("build AST finished")
	return output[0], nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1

	case "*", "/":
		return 2

	case "^":
		return 3

	default:
		return 0
	}
}
