package user

import "time"

// CreateDTO is a creation user dto.
type CreateDTO struct {
	Email    string
	Username string
	Password string
}

// UpdateDTO is an update user dto.
type UpdateDTO struct {
	Username  *string
	Email     *string
	Bio       *string
	Image     *string
	UpdatedAt time.Time
}
