// Package domain provides a domain service.
package domain

import (
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/maypok86/conduit/internal/repository/psql"
)

// Services is a collection of all services in the system.
type Services struct {
	User    user.Service
	Profile profile.Service
}

// NewServices returns a new instance of Services.
func NewServices(repositories psql.Repositories, passwordHasher user.PasswordHasher) Services {
	return Services{
		User:    user.NewService(repositories.User, passwordHasher),
		Profile: profile.NewService(repositories.Profile),
	}
}
