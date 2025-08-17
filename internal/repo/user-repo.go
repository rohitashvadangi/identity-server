package repo

import "github.com/rohitashvadangi/identity-server/internal/proto"

type UserRepository interface {
	// Find a user by username
	FindByUsername(username string) (*proto.User, error)

	// Save a new user
	Save(user *proto.User) error
}
