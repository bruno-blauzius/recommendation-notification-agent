package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/project-go-sender-recommendation-agent/internal/infrastructure/config"
)

func NewRedisConnection(cfg *config.Config) (*redis.Client, error) {

	if cfg == nil {
		return nil, fmt.Errorf("database.newRedisConnection: config is nil")
	}

	url := cfg.RedisAddr()
	password := cfg.RedisPassword
	username := cfg.RedisUsername

	rdb := redis.NewClient(&redis.Options{
		Addr:            url,
		Password:        password,
		Username:        username,
		ReadBufferSize:  1024 * 1024,
		WriteBufferSize: 1024 * 1024,
		Protocol:        3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		if strings.Contains(err.Error(), "WRONGPASS") || strings.Contains(err.Error(), "NOAUTH") {
			return nil, fmt.Errorf("database.newRedisConnection: authentication failed for Redis at %s", url)
		}
		return nil, fmt.Errorf("database.newRedisConnection ping: %w", err)
	}

	return rdb, nil
}
