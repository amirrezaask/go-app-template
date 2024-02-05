package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	conn *redis.Client
}

func NewRedis(ctx context.Context, hostPort string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: hostPort,
		DB:   db,
	})
	statusCmd := client.Ping(ctx)
	if err := statusCmd.Err(); err != nil {
		return nil, err
	}
	return &Redis{
		conn: client,
	}, nil
}
