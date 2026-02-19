package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     int    `env:"REDIS_PORT,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type redisConfig struct {
	raw redisEnvConfig
}

// NewRedisConfig creates a new Redis configuration from environment variables.
func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &redisConfig{raw: raw}, nil
}

// Host returns the Redis host.
func (c *redisConfig) Host() string {
	return c.raw.Host
}

// Port returns the Redis port.
func (c *redisConfig) Port() int {
	return c.raw.Port
}

// Password returns the Redis password.
func (c *redisConfig) Password() string {
	return c.raw.Password
}

// DB returns the Redis database number.
func (c *redisConfig) DB() int {
	return c.raw.DB
}

// Address returns the Redis address in host:port format.
func (c *redisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.raw.Host, c.raw.Port)
}
