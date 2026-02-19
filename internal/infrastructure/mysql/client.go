package mysql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config represents MySQL configuration.
type Config interface {
	DSN() string
	MaxOpenConns() int
	MaxIdleConns() int
	ConnMaxLifetime() time.Duration
}

// Client wraps the database connection pool.
type Client struct {
	db *sql.DB
}

// NewClient creates a new MySQL client with connection pooling.
func NewClient(cfg Config) (*Client, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns())
	db.SetMaxIdleConns(cfg.MaxIdleConns())
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime())

	return &Client{db: db}, nil
}

// Ping verifies the database connection.
func (c *Client) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Close closes the database connection.
func (c *Client) Close() error {
	return c.db.Close()
}

// DB returns the underlying database connection.
func (c *Client) DB() *sql.DB {
	return c.db
}

// BeginTx starts a new transaction.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}
