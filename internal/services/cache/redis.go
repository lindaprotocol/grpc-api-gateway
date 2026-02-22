package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
}

// NewRedisClientFromConfig creates a Redis client from application config
func NewRedisClientFromConfig(cfg config.RedisConfig) (*RedisClient, error) {
	return NewRedisClient(RedisConfig{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
	})
}

func NewRedisClient(cfg RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis
func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, ttl).Err()
}

// Get retrieves a value from Redis
func (r *RedisClient) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(key string) bool {
	result, err := r.client.Exists(r.ctx, key).Result()
	return err == nil && result > 0
}

// Incr increments a counter
func (r *RedisClient) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Expire sets expiration on a key
func (r *RedisClient) Expire(key string, ttl time.Duration) error {
	return r.client.Expire(r.ctx, key, ttl).Err()
}

// GetTTL returns the remaining TTL of a key
func (r *RedisClient) GetTTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Client returns the underlying Redis client for advanced operations (e.g. rate limiting)
func (r *RedisClient) Client() *redis.Client {
	return r.client
}