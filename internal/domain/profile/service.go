package profile

import (
	"context"
	"fmt"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=profile_test

// Repository is a profile repository.
type Repository interface {
	GetByUsername(ctx context.Context, username string) (Profile, error)
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
