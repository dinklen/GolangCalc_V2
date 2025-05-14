package dberr

import "errors"

var (
	// Database errors
	ErrDatabaseCreatedFailed = errors.New("failed to create the database")
	ErrDatabasePingFailed    = errors.New("failed to ping the database")

	// Migrations errors
	ErrMigrationsReadingFailed = errors.New("failed to read migrations")
	ErrMigrationsUpFailed      = errors.New("migrations up failed")

	// Accounts requests errors
	ErrAccountGettingFailed  = errors.New("failed to get the account")
	ErrAccountCreatingFailed = errors.New("failed to create the acount")

	// Expressions requests errors
	ErrExpressionCreatingFailed      = errors.New("failed to create the expression")
	ErrExpressionUpdatingStateFailed = errors.New("failed to update the expression state")
	ErrExpressionsGettingFailed      = errors.New("failed to get the expressions") // Multiple expressions
	ErrExpressionsParsingFailed      = errors.New("failed to parse the expressions")
	ErrExpressionGettingFailed       = errors.New("failed to get the expression") // Current expression

	// Subexpressions requests errors
	ErrParentExpressionIDParsingFailed  = errors.New("failed to parse the parent expression ID")
	ErrSubexpressionCreatingFailed      = errors.New("failed to create the subexpression")
	ErrSubexpressionUpdatingStateFailed = errors.New("failed to update the subexpression state")
)
