package parser

import (
	"strings"

	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calcerr"

	"go.uber.org/zap"
)

func Tokenize(expr string, logger *zap.Logger) ([]Token, error) {
	logger.Info("tokenize started")

	var tokens []Token
	var buffer strings.Builder

	for i, ch := range expr {
		switch {
		case ch == ' ' || ch == '\t':
			continue

		case ch >= '0' && ch <= '9' || ch == '.':
			buffer.WriteRune(ch)

		default:
			if buffer.Len() > 0 {
				tokens = append(tokens, Token{Type: Number, Value: buffer.String()})

				buffer.Reset()
			}

			switch ch {
			case '+', '-', '*', '/', '^':
				if ch == '-' && (i == 0 || tokens[len(tokens)-1].Type == LeftParen || (tokens[len(tokens)-1].Type == Operator && tokens[len(tokens)-1].Value != ")")) {
					tokens = append(tokens, Token{Type: Number, Value: "0"})
				}

				tokens = append(tokens, Token{Type: Operator, Value: string(ch)})

			case '(':
				tokens = append(tokens, Token{Type: LeftParen, Value: string(ch)})

			case ')':
				tokens = append(tokens, Token{Type: RightParen, Value: string(ch)})

			default:
				return nil, calcerr.ErrUnknownCharacter
			}
		}
	}

	if buffer.Len() > 0 {
		tokens = append(tokens, Token{Type: Number, Value: buffer.String()})
	}

	logger.Info("tokenize finished")
	return tokens, nil
}
