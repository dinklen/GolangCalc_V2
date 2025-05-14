package neterr

import "errors"

var (
	// HTTPS errors
	ErrShuttingDownFailed = errors.New("failed to shutdown the server")
	ErrStartingFailed     = errors.New("faield to start the server")

	// gRPC errors
	ErrGRPCServerStartingFailed   = errors.New("failed to start the gRPC server")
	ErrGRPCServerListeningFailed  = errors.New("gRPC server listening failed")
	ErrGRPCServerServingFailed    = errors.New("gRPC server serving failed")
	ErrGRPCServerConnectionClosed = errors.New("client connection closed")
)
