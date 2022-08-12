// Package repository represents a repository layer.
package repository

import "github.com/maypok86/conduit/pkg/postgres"

// UserPostgres is a user repository.
type UserPostgres struct {
	db *postgres.Postgres
}

// NewUserPostgres creates a new UserPostgres.
func NewUserPostgres(db *postgres.Postgres) UserPostgres {
	return UserPostgres{
		db: db,
	}
}
