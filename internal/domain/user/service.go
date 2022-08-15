package user

import (
	"context"
	"fmt"
	"time"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=user_test

// Repository is a user repository.
type Repository interface {
	Create(ctx context.Context, dto User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	UpdateByEmail(ctx context.Context, email string, updateDTO UpdateDTO) (User, error)
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

// NewService creates a new user service.
func NewService(userRepository Repository, passwordHasher PasswordHasher) Service {
	return Service{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

// Create creates a new user.
func (s Service) Create(ctx context.Context, dto CreateDTO) (User, error) {
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

	user, err = s.userRepository.Create(ctx, user)
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

// Login provides user login.
func (s Service) Login(ctx context.Context, email, password string) (User, error) {
	user, err := s.GetByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	if err := s.passwordHasher.Check(password, user.Password); err != nil {
		return User{}, fmt.Errorf("can not check password: %w", err)
	}

	return user, nil
}

// UpdateByEmail updates user by email.
func (s Service) UpdateByEmail(ctx context.Context, email string, dto UpdateDTO) (User, error) {
	dto.UpdatedAt = time.Now()

	user, err := s.userRepository.UpdateByEmail(ctx, email, dto)
	if err != nil {
		return User{}, fmt.Errorf("can not update user: %w", err)
	}

	return user, nil
}
