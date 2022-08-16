// Package profile represents a profile domain.
package profile

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrNotFound is an error that indicates that profile not found.
var ErrNotFound = errors.New("profile not found")

// Profile is a profile entity.
type Profile struct {
	ID        uuid.UUID
	Username  string
	Bio       *string
	Image     *string
	Following bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetBio returns bio.
func (p Profile) GetBio() string {
	if p.Bio == nil {
		return ""
	}

	return *p.Bio
}

// GetImage returns image.
func (p Profile) GetImage() string {
	if p.Image == nil {
		return ""
	}

	return *p.Image
}
