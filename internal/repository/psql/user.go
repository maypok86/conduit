// Package psql represents a repository for PostgreSQL.
package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
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

// Create creates a new user.
func (ur UserRepository) Create(ctx context.Context, dto user.User) (user.User, error) {
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

// GetByEmail returns user by email.
func (ur UserRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	sql, args, err := ur.db.Builder.Select(
		"id",
		"username",
		"password",
		"bio",
		"image",
		"created_at",
		"updated_at",
	).From("users").Where("email = ?", email).ToSql()
	if err != nil {
		return user.User{}, fmt.Errorf("can not build select user query: %w", err)
	}

	u := user.User{Email: email}
	if err := ur.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Bio,
		&u.Image,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, fmt.Errorf("can not find user: %w", user.ErrNotFound)
		}

		return user.User{}, fmt.Errorf("can not find user: %w", err)
	}

	return u, nil
}
