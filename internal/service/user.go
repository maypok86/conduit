// Package service represents a service layer.
package service

// UserRepository is a user repository interface.
type UserRepository interface{}

// UserService is a user service interface.
type UserService struct {
	userRepository UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepository UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}
