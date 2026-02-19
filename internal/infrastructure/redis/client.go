package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config represents Redis configuration.
type Config interface {
	Address() string
	Password() string
	DB() int
}

// Client wraps the Redis client.
type Client struct {
	rdb *redis.Client
}

// NewClient creates a new Redis client.
func NewClient(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password(),
		DB:       cfg.DB(),
	})

	return &Client{rdb: rdb}, nil
}

// Ping verifies the Redis connection.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// Client returns the underlying Redis client.
func (c *Client) Client() *redis.Client {
	return c.rdb
}

// Set stores a value with the given key and TTL.
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value by key.
func (c *Client) Get(ctx context.Context, key string, dest any) error {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key.
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Exists checks if a key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}
