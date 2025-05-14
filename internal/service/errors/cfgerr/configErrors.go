package cfgerr

import "errors"

var (
	ErrConfigReadingFailed          = errors.New("failed to read the config file")
	ErrConfigDataUnmarshalingFailed = errors.New("failed to unmarshal the config data")
	ErrConfigUpdattingFailed        = errors.New("failed to update the config")
)
