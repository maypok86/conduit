package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/maypok86/conduit/pkg/postgres"
	"go.uber.org/zap"
)

// ProfileRepository is a profile repository.
type ProfileRepository struct {
	db *postgres.Postgres
}

// NewProfileRepository creates a new ProfileRepository.
func NewProfileRepository(db *postgres.Postgres) ProfileRepository {
	return ProfileRepository{
		db: db,
	}
}

// GetByUsername returns profile by username.
func (pr ProfileRepository) GetByUsername(ctx context.Context, username string) (profile.Profile, error) {
	sql, args, err := pr.db.Builder.Select(
		"id",
		"bio",
		"image",
		"created_at",
		"updated_at",
	).From("users").Where(sq.Eq{"username": username}).Limit(1).ToSql()
	if err != nil {
		return profile.Profile{}, fmt.Errorf("can not build select profile by username query: %w", err)
	}

	logger.FromContext(ctx).Debug("select profile by username query", zap.String("sql", sql), zap.Any("args", args))

	p := profile.Profile{Username: username}
	if err := pr.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&p.ID,
		&p.Bio,
		&p.Image,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return profile.Profile{}, fmt.Errorf("can not find profile by username: %w", profile.ErrNotFound)
		}

		return profile.Profile{}, fmt.Errorf("can not find profile by username: %w", err)
	}

	return p, nil
}

// GetByEmail returns profile by username.
func (pr ProfileRepository) GetByEmail(ctx context.Context, email string) (profile.Profile, error) {
	sql, args, err := pr.db.Builder.Select(
		"id",
		"username",
		"bio",
		"image",
		"created_at",
		"updated_at",
	).From("users").Where(sq.Eq{"email": email}).Limit(1).ToSql()
	if err != nil {
		return profile.Profile{}, fmt.Errorf("can not build select profile by email query: %w", err)
	}

	logger.FromContext(ctx).Debug("select profile by username query", zap.String("sql", sql), zap.Any("args", args))

	var p profile.Profile
	if err := pr.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&p.ID,
		&p.Username,
		&p.Bio,
		&p.Image,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return profile.Profile{}, fmt.Errorf("can not find profile by email: %w", profile.ErrNotFound)
		}

		return profile.Profile{}, fmt.Errorf("can not find profile by email: %w", err)
	}

	return p, nil
}

// CheckFollowing returns profile by id with follow.
func (pr ProfileRepository) CheckFollowing(
	ctx context.Context,
	email string,
	followeeID uuid.UUID,
) (profile.Profile, error) {
	sql, args, err := pr.db.Builder.Select(
		"id",
		"username",
		"bio",
		"image",
		"users.created_at",
		"users.updated_at",
	).From("users").Join("follows ON users.id = follows.follower_id").Where(
		sq.And{sq.Eq{"email": email}, sq.Eq{"followee_id": followeeID}},
	).Limit(1).ToSql()
	if err != nil {
		return profile.Profile{}, fmt.Errorf("can not build check following profile query: %w", err)
	}

	logger.FromContext(ctx).Debug("select profile with follow query", zap.String("sql", sql), zap.Any("args", args))

	p := profile.Profile{
		Following: true,
	}
	if err := pr.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&p.ID,
		&p.Username,
		&p.Bio,
		&p.Image,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return profile.Profile{}, fmt.Errorf("can not check following profile: %w", profile.ErrNotFound)
		}

		return profile.Profile{}, fmt.Errorf("can not check following profile: %w", err)
	}

	return p, nil
}

// Follow adds follow relationship.
func (pr ProfileRepository) Follow(ctx context.Context, followeeID, followerID uuid.UUID) error {
	sql, args, err := pr.db.Builder.Insert("follows").
		Columns("followee_id", "follower_id").
		Values(followeeID, followerID).
		ToSql()
	if err != nil {
		return fmt.Errorf("can not build follow query: %w", err)
	}

	logger.FromContext(ctx).Debug("follow query", zap.String("sql", sql), zap.Any("args", args))

	if _, err := pr.db.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("can not follow: %w", err)
	}

	return nil
}

// Unfollow removes follow relationship.
func (pr ProfileRepository) Unfollow(ctx context.Context, followeeID, followerID uuid.UUID) error {
	sql, args, err := pr.db.Builder.Delete("follows").
		Where(sq.And{sq.Eq{"followee_id": followeeID}, sq.Eq{"follower_id": followerID}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("can not build unfollow query: %w", err)
	}

	logger.FromContext(ctx).Debug("unfollow query", zap.String("sql", sql), zap.Any("args", args))

	if _, err := pr.db.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("can not unfollow: %w", err)
	}

	return nil
}
