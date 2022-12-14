// Package psql represents a repository for PostgreSQL.
package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/maypok86/conduit/pkg/postgres"
	"go.uber.org/zap"
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

	logger.FromContext(ctx).Debug("create user query", zap.String("sql", sql), zap.Any("args", args))

	if err := ur.db.Pool.QueryRow(ctx, sql, args...).Scan(&dto.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return user.User{}, fmt.Errorf("can not insert user: %w", user.ErrAlreadyExist)
			}
		}

		return user.User{}, fmt.Errorf("can not insert user: %w", err)
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
	).From("users").Where(sq.Eq{"email": email}).Limit(1).ToSql()
	if err != nil {
		return user.User{}, fmt.Errorf("can not build select user by email query: %w", err)
	}

	logger.FromContext(ctx).Debug("select user by email query", zap.String("sql", sql), zap.Any("args", args))

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
			return user.User{}, fmt.Errorf("can not find user by email: %w", user.ErrNotFound)
		}

		return user.User{}, fmt.Errorf("can not find user by email: %w", err)
	}

	return u, nil
}

func (ur UserRepository) buildUpdateUserQuery(updateBuilder sq.UpdateBuilder, dto user.UpdateDTO) sq.UpdateBuilder {
	if dto.Username != nil {
		updateBuilder = updateBuilder.Set("username", *dto.Username)
	}

	if dto.Email != nil {
		updateBuilder = updateBuilder.Set("email", *dto.Email)
	}

	if dto.Bio != nil {
		updateBuilder = updateBuilder.Set("bio", *dto.Bio)
	}

	if dto.Image != nil {
		updateBuilder = updateBuilder.Set("image", *dto.Image)
	}

	return updateBuilder.Set("updated_at", dto.UpdatedAt)
}

// UpdateByEmail updates user by email.
func (ur UserRepository) UpdateByEmail(ctx context.Context, email string, dto user.UpdateDTO) (user.User, error) {
	updateBuilder := ur.buildUpdateUserQuery(ur.db.Builder.Update("users"), dto)

	sql, args, err := updateBuilder.Suffix(
		"RETURNING id, username, email, password, bio, image, created_at",
	).Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return user.User{}, fmt.Errorf("can not build update user by email query: %w", err)
	}

	logger.FromContext(ctx).Debug("update user by email query", zap.String("sql", sql), zap.Any("args", args))

	u := user.User{UpdatedAt: dto.UpdatedAt}
	if err := ur.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Bio,
		&u.Image,
		&u.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, fmt.Errorf("can not update user by email: %w", user.ErrNotFound)
		}

		return user.User{}, fmt.Errorf("can not update user by email: %w", err)
	}

	return u, nil
}
