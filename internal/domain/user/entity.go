// Package user represents a user domain.
package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrAlreadyExist is an error that indicates that user already exists.
var ErrAlreadyExist = errors.New("user with given email or nickname already exist")

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

func (u User) GetBio() string {
	if u.Bio == nil {
		return ""
	}

	return *u.Bio
}

func (u User) GetImage() string {
	if u.Image == nil {
		return ""
	}

	return *u.Image
}
