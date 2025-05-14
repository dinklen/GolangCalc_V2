package tmerr

import "errors"

var (
	// Stream error
	ErrStreamCreatingFailed = errors.New("the stream creating failed")

	// Client error
	ErrClientCreatingFailed    = errors.New("failed to create the gRPC client")
	ErrClientConnectionRefused = errors.New("client connection refused")

	// Evaluating error
	ErrTimeout = errors.New("timeout")
)
