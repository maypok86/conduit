// Package user represents a user domain.
package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrAlreadyExist is an error that indicates that user already exists.
	ErrAlreadyExist = errors.New("user with given email or nickname already exist")
	// ErrNotFound is an error that indicates that user not found.
	ErrNotFound = errors.New("user not found")
)

// User is a user entity.
type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	Bio       *string
	Image     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetBio returns bio.
func (u User) GetBio() string {
	if u.Bio == nil {
		return ""
	}

	return *u.Bio
}

// GetImage returns image.
func (u User) GetImage() string {
	if u.Image == nil {
		return ""
	}

	return *u.Image
}
