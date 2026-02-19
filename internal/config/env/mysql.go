package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type mysqlEnvConfig struct {
	Host            string `env:"MYSQL_HOST,required"`
	Port            int    `env:"MYSQL_PORT,required"`
	User            string `env:"MYSQL_USER,required"`
	Password        string `env:"MYSQL_PASSWORD,required"`
	Database        string `env:"MYSQL_DATABASE,required"`
	MaxOpenConns    int    `env:"MYSQL_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int    `env:"MYSQL_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime string `env:"MYSQL_CONN_MAX_LIFETIME" envDefault:"5m"`
}

type mysqlConfig struct {
	raw mysqlEnvConfig
}

// NewMySQLConfig creates a new MySQL configuration from environment variables.
func NewMySQLConfig() (*mysqlConfig, error) {
	var raw mysqlEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &mysqlConfig{raw: raw}, nil
}

// Host returns the MySQL host.
func (c *mysqlConfig) Host() string {
	return c.raw.Host
}

// Port returns the MySQL port.
func (c *mysqlConfig) Port() int {
	return c.raw.Port
}

// User returns the MySQL user.
func (c *mysqlConfig) User() string {
	return c.raw.User
}

// Password returns the MySQL password.
func (c *mysqlConfig) Password() string {
	return c.raw.Password
}

// Database returns the MySQL database name.
func (c *mysqlConfig) Database() string {
	return c.raw.Database
}

// MaxOpenConns returns the maximum number of open connections.
func (c *mysqlConfig) MaxOpenConns() int {
	return c.raw.MaxOpenConns
}

// MaxIdleConns returns the maximum number of idle connections.
func (c *mysqlConfig) MaxIdleConns() int {
	return c.raw.MaxIdleConns
}

// ConnMaxLifetime returns the maximum connection lifetime.
func (c *mysqlConfig) ConnMaxLifetime() time.Duration {
	d, err := time.ParseDuration(c.raw.ConnMaxLifetime)
	if err != nil {
		return 5 * time.Minute
	}
	return d
}

// DSN returns the MySQL connection string.
func (c *mysqlConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		c.raw.User, c.raw.Password, c.raw.Host, c.raw.Port, c.raw.Database)
}
