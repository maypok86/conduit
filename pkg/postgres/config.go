package postgres

import (
	"fmt"
	"time"
)

// ConnectionConfig is a configuration for postgres connection.
type ConnectionConfig struct {
	host     string
	port     string
	dbname   string
	username string
	password string
	sslmode  string
}

// NewConnectionConfig creates a new ConnectionConfig.
func NewConnectionConfig(host, port, dbname, username, password, sslmode string) ConnectionConfig {
	return ConnectionConfig{
		host:     host,
		port:     port,
		dbname:   dbname,
		username: username,
		password: password,
		sslmode:  sslmode,
	}
}

func (cc ConnectionConfig) getDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cc.username,
		cc.password,
		cc.host,
		cc.port,
		cc.dbname,
		cc.sslmode,
	)
}

// Option is a functional option for configuring a Postgres.
type Option func(*Postgres)

// WithMaxPoolSize sets the max pool size for the Postgres.
func WithMaxPoolSize(maxPoolSize int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = maxPoolSize
	}
}

// WithConnAttempts sets the max attempts for connecting to the Postgres.
func WithConnAttempts(connAttempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = connAttempts
	}
}

// WithConnTimeout sets the timeout for connecting to the Postgres.
func WithConnTimeout(connTimeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = connTimeout
	}
}
