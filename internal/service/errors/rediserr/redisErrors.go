package rediserr

import "errors"

var ErrPingFailed = errors.New("failed to ping Redis server")
