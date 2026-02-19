package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type jwtEnvConfig struct {
	AccessTokenTTL string `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"15m"`
	Issuer         string `env:"JWT_ISSUER" envDefault:"tasks-service"`
	SecretKey      string `env:"JWT_SECRET_KEY,required"`
}

type jwtConfig struct {
	raw jwtEnvConfig
}

// NewJWTConfig creates a new JWT configuration from environment variables.
func NewJWTConfig() (*jwtConfig, error) {
	var raw jwtEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &jwtConfig{raw: raw}, nil
}

// AccessTokenTTL returns the access token time-to-live duration.
func (c *jwtConfig) AccessTokenTTL() time.Duration {
	ttl, err := time.ParseDuration(c.raw.AccessTokenTTL)
	if err != nil {
		return 15 * time.Minute
	}
	return ttl
}

// Issuer returns the JWT issuer claim value.
func (c *jwtConfig) Issuer() string {
	return c.raw.Issuer
}

// SecretKey returns the JWT secret key.
func (c *jwtConfig) SecretKey() string {
	return c.raw.SecretKey
}
