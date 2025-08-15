package redis_client

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string, pass string, num int, protocol int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       num,
		Protocol: protocol,
	})

	return &RedisClient{client: client}
}

// QueueClient interface'ini implement eder
func (r *RedisClient) Peek(ctx context.Context, queueName string) bool {
	return false
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
