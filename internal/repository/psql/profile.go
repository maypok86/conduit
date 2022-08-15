package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/maypok86/conduit/pkg/postgres"
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
		"bio",
		"image",
		"created_at",
		"updated_at",
	).From("users").Where(sq.Eq{"username": username}).Limit(1).ToSql()
	if err != nil {
		return profile.Profile{}, fmt.Errorf("can not build select profile by username query: %w", err)
	}

	p := profile.Profile{Username: username}
	if err := pr.db.Pool.QueryRow(ctx, sql, args...).Scan(
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
