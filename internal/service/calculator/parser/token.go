package parser

type TokenType int

const (
	Number TokenType = iota
	Operator
	LeftParen
	RightParen
)

type Token struct {
	Type  TokenType
	Value string
}
