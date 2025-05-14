package calcerr

import "errors"

var (
	// Syntax errors
	ErrNotEnoughOperands     = errors.New("not enough operands")
	ErrMismatchedParentheses = errors.New("mismatched parentheses")
	ErrInvalidExpression     = errors.New("invalid expression")
	ErrUnknownCharacter      = errors.New("unknown character")

	// Evaluation errors
	ErrEmtpyNode                 = errors.New("empty node ecountered")
	ErrUnexpectedCharacter       = errors.New("unexpected character")
	ErrBadAnswerFromMicroservice = errors.New("bad answer from microservice")
	ErrDivisionByZero            = errors.New("division by zero")
	ErrUnknownOperator           = errors.New("unknown operator")
	ErrParsingFailed             = errors.New("parsing failed")
)
