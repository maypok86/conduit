// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxPoolSize  int           = 1
	defaultConnAttempts int           = 10
	defaultConnTimeout  time.Duration = time.Second
)

// ErrUnableToConnect is returned when unable to connect to the postgres.
var ErrUnableToConnect = errors.New("all attempts are exceeded. Unable to connect to instance")

// Postgres is a postgres connection.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	Builder      sq.StatementBuilderType
	Pool         *pgxpool.Pool
}

// New creates a new postgres connection.
func New(ctx context.Context, connectionConfig ConnectionConfig, opts ...Option) (*Postgres, error) {
	instance := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(instance)
	}

	instance.Builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	dsn := connectionConfig.getDSN()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	poolCfg.MaxConns = int32(instance.maxPoolSize)

	for instance.connAttempts > 0 {
		instance.Pool, err = pgxpool.ConnectConfig(ctx, poolCfg)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", instance.connAttempts)
		time.Sleep(instance.connTimeout)

		instance.connAttempts--
	}

	if err != nil {
		return nil, ErrUnableToConnect
	}

	return instance, nil
}

// Close closes the postgres connection.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
