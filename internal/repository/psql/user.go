// Package psql represents a repository for PostgreSQL.
package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/maypok86/conduit/pkg/postgres"
)

// UserRepository is a user repository.
type UserRepository struct {
	db *postgres.Postgres
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *postgres.Postgres) UserRepository {
	return UserRepository{
		db: db,
	}
}

// CreateUser creates a new user.
func (ur UserRepository) CreateUser(ctx context.Context, dto user.User) (user.User, error) {
	sql, args, err := ur.db.Builder.Insert("users").Columns(
		"email",
		"username",
		"password",
	).Suffix("RETURNING id").Values(
		dto.Email,
		dto.Username,
		dto.Password,
	).ToSql()
	if err != nil {
		return user.User{}, fmt.Errorf("can not build insert user query: %w", err)
	}

	if err := ur.db.Pool.QueryRow(ctx, sql, args...).Scan(&dto.ID); err != nil {
		const errorFmtString = "can not insert user: %w"

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return user.User{}, fmt.Errorf(errorFmtString, user.ErrAlreadyExist)
			}
		}

		return user.User{}, fmt.Errorf(errorFmtString, err)
	}

	return dto, nil
}
