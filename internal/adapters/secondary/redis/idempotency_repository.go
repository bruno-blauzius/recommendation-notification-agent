package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type idempotencyRepository struct {
	client *redis.Client
}

// NewIdempotencyRepository creates a new Redis-backed IdempotencyRepository.
func NewIdempotencyRepository(client *redis.Client) *idempotencyRepository {
	return &idempotencyRepository{client: client}
}

// Set stores key with value and the given TTL.
// Typically called with a "processing:{id}" key when a record enters processing.
func (r *idempotencyRepository) Set(key string, value string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("idempotencyRepository.Set key=%s: %w", key, err)
	}
	return nil
}

// Get retrieves the value for key.
// Returns ("", nil) when the key does not exist.
func (r *idempotencyRepository) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("idempotencyRepository.Get key=%s: %w", key, err)
	}
	return value, nil
}

// Exists reports whether key is present in Redis.
func (r *idempotencyRepository) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("idempotencyRepository.Exists key=%s: %w", key, err)
	}
	return count > 0, nil
}

func (r *idempotencyRepository) GenerateKey(prefix string, value string) string {
	return fmt.Sprintf("%s:%s", prefix, value)
}
