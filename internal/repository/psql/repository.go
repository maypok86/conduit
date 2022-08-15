package psql

import "github.com/maypok86/conduit/pkg/postgres"

// Repositories is a collection of all repositories in the system.
type Repositories struct {
	User    UserRepository
	Profile ProfileRepository
}

// NewRepositories returns a new instance of Repositories.
func NewRepositories(db *postgres.Postgres) Repositories {
	return Repositories{
		User:    NewUserRepository(db),
		Profile: NewProfileRepository(db),
	}
}
