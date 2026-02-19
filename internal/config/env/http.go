package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	Host        string `env:"HTTP_HOST,required"`
	Port        string `env:"HTTP_PORT,required"`
	ReadTimeout string `env:"HTTP_READ_TIMEOUT,required"`
}

type httpConfig struct {
	raw httpEnvConfig
}

// NewHTTPConfig creates a new HTTP configuration from environment variables.
func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &httpConfig{raw: raw}, nil
}

// Host returns the HTTP server host.
func (h *httpConfig) Host() string {
	return h.raw.Host
}

// Port returns the HTTP server port.
func (h *httpConfig) Port() string {
	return h.raw.Port
}

// Address returns the full HTTP server address.
func (h *httpConfig) Address() string {
	return fmt.Sprintf("%s:%s", h.raw.Host, h.raw.Port)
}

// ReadTimeout returns the HTTP read timeout duration.
func (h *httpConfig) ReadTimeout() time.Duration {
	readTimeout, err := time.ParseDuration(h.raw.ReadTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return readTimeout
}
