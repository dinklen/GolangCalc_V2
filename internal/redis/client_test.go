package redisc

import (
	"context"
	"testing"
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/rediserr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewClient_Success(t *testing.T) {
	cfg := &config.Config{}
	cfg.Redis.Host = "localhost"
	cfg.Redis.Port = "6379"
	cfg.Redis.Password = ""
	cfg.Redis.WaitingTime = time.Second

	logger := zaptest.NewLogger(t)
	client, err := NewClient(cfg, logger)

	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx := context.Background()
	err = client.Ping(ctx).Err()
	assert.NoError(t, err)
}

func TestNewClient_InvalidAddress(t *testing.T) {
	cfg := &config.Config{}
	cfg.Redis.Host = "invalid_host"
	cfg.Redis.Port = "6379"
	cfg.Redis.Password = ""
	cfg.Redis.WaitingTime = time.Nanosecond

	logger := zaptest.NewLogger(t)
	client, err := NewClient(cfg, logger)

	assert.ErrorIs(t, err, rediserr.ErrPingFailed)
	assert.Nil(t, client)
}

func TestNewClient_ContextTimeout(t *testing.T) {
	cfg := &config.Config{}
	cfg.Redis.Host = "localhost"
	cfg.Redis.Port = "6379"
	cfg.Redis.Password = ""
	cfg.Redis.WaitingTime = time.Nanosecond

	logger := zaptest.NewLogger(t)
	client, err := NewClient(cfg, logger)

	assert.ErrorIs(t, err, rediserr.ErrPingFailed)
	assert.Nil(t, client)
}
