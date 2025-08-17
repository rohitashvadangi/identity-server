package stores

import (
	"github.com/rohitashvadangi/identity-server/internal/proto"
)

type TokenStore interface {
	// Save token with associated user ID
	Save(token *proto.Token)

	// Retrieve user ID by token
	Get(token string) (*proto.Token, error)

	// Delete token
	Delete(token string)
}
