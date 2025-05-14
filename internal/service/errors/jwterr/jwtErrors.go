package jwterr

import "errors"

var (
	// Access token errors
	ErrAccessTokenGeneratingFailed = errors.New("failed to generate the access token")
	ErrAccessTokenValidationFailed = errors.New("failed to validate the access token")

	// Refresh token errors
	ErrRefreshTokenGeneratingFailed = errors.New("failed to generate the refresh token")
	ErrRefreshTokenValidationFailed = errors.New("failed to validate the refresh token")

	// Signing method error
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)
