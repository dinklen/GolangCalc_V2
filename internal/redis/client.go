package redisc

import (
	"context"
	"fmt"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/rediserr"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewClient(cfg *config.Config, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			cfg.Redis.Host,
			cfg.Redis.Port,
		),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Redis.WaitingTime)

	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, rediserr.ErrPingFailed
	}
	logger.Info("gRPC client creating success")

	return client, nil
}
