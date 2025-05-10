package redis

import (
	"cobtext"
	"time"
	"fmt"

	"github.con/dinklen/GolangCalc_V2/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			cfg.Redis.Host,
			cfg.Redis.Port,
		),
		Password: cfg.Redis.Password,
		DB: 0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Redis.WaitingTime)

	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	
	return client, nil
}
