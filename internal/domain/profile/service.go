package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=profile_test

// Repository is a profile repository.
type Repository interface {
	GetByUsername(ctx context.Context, username string) (Profile, error)
	GetByEmail(ctx context.Context, email string) (Profile, error)
	CheckFollowing(ctx context.Context, followeeID, followerID uuid.UUID) error
	Follow(ctx context.Context, followeeID, followerID uuid.UUID) error
	Unfollow(ctx context.Context, followeeID, followerID uuid.UUID) error
}

// Service a profile service interface.
type Service struct {
	profileRepository Repository
}

// NewService creates a new profile service.
func NewService(profileRepository Repository) Service {
	return Service{
		profileRepository: profileRepository,
	}
}

// GetByUsername gets a profile by username.
func (s Service) GetByUsername(ctx context.Context, username string) (Profile, error) {
	profile, err := s.profileRepository.GetByUsername(ctx, username)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to get profile by username: %w", err)
	}

	return profile, nil
}

// GetByEmail gets a profile by email.
func (s Service) GetByEmail(ctx context.Context, email string) (Profile, error) {
	profile, err := s.profileRepository.GetByEmail(ctx, email)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to get profile by email: %w", err)
	}

	return profile, nil
}

// GetWithFollow gets a profile with follow checking.
func (s Service) GetWithFollow(ctx context.Context, email, username string) (Profile, error) {
	followee, err := s.GetByUsername(ctx, username)
	if err != nil {
		return Profile{}, err
	}

	follower, err := s.GetByEmail(ctx, email)
	if err != nil {
		return Profile{}, err
	}

	if err := s.profileRepository.CheckFollowing(ctx, followee.ID, follower.ID); err != nil {
		if errors.Is(err, ErrNotFound) {
			followee.Following = false
			return followee, nil
		}

		return Profile{}, fmt.Errorf("failed to check following: %w", err)
	}

	followee.Following = true

	return followee, nil
}

// Follow make a follow relationship.
func (s Service) Follow(ctx context.Context, email, username string) (Profile, error) {
	followee, err := s.GetByUsername(ctx, username)
	if err != nil {
		return Profile{}, err
	}

	follower, err := s.GetByEmail(ctx, email)
	if err != nil {
		return Profile{}, err
	}

	if err := s.profileRepository.Follow(ctx, followee.ID, follower.ID); err != nil {
		return Profile{}, fmt.Errorf("failed to follow: %w", err)
	}

	followee.Following = true

	return followee, nil
}

// Unfollow delete a follow relationship.
func (s Service) Unfollow(ctx context.Context, email, username string) (Profile, error) {
	followee, err := s.GetByUsername(ctx, username)
	if err != nil {
		return Profile{}, err
	}

	follower, err := s.GetByEmail(ctx, email)
	if err != nil {
		return Profile{}, err
	}

	if err := s.profileRepository.Unfollow(ctx, followee.ID, follower.ID); err != nil {
		return Profile{}, fmt.Errorf("failed to unfollow: %w", err)
	}

	followee.Following = false

	return followee, nil
}
