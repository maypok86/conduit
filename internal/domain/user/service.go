package user

import (
	"context"
	"fmt"
	"time"
)

// Repository is a user repository.
type Repository interface {
	CreateUser(ctx context.Context, dto User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

// PasswordHasher is a password hasher.
type PasswordHasher interface {
	Hash(string) (string, error)
	Check(string, string) error
}

// Service is a user service interface.
type Service struct {
	userRepository Repository
	passwordHasher PasswordHasher
}

// NewService creates a new UserService.
func NewService(userRepository Repository, passwordHasher PasswordHasher) Service {
	return Service{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

// CreateUser creates a new user.
func (s Service) CreateUser(ctx context.Context, dto CreateDTO) (User, error) {
	passwordHash, err := s.passwordHasher.Hash(dto.Password)
	if err != nil {
		return User{}, fmt.Errorf("can not hash password: %w", err)
	}

	now := time.Now()
	user := User{
		Email:     dto.Email,
		Username:  dto.Username,
		Password:  passwordHash,
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return User{}, fmt.Errorf("can not create user: %w", err)
	}

	return user, nil
}

// GetByEmail returns user by email.
func (s Service) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("can not get user by email: %w", err)
	}

	return user, nil
}
