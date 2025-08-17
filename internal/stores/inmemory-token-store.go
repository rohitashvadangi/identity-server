package stores

import (
	"errors"
	"sync"
	"time"

	"github.com/rohitashvadangi/identity-server/internal/proto"
)

type InMemoryTokenStore struct {
	tokens map[string]*proto.Token
	mu     sync.RWMutex
}

// Ensure InMemoryTokenStore implements TokenStore
var _ TokenStore = (*InMemoryTokenStore)(nil)

func NewTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		tokens: make(map[string]*proto.Token),
	}
}

func (s *InMemoryTokenStore) Save(token *proto.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token.TokenID] = token
}

func (s *InMemoryTokenStore) Get(tokenID string) (*proto.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tokens[tokenID]
	if !ok || time.Now().After(t.ExpiresAt) {
		return nil, errors.New("token not found or expired")
	}
	return t, nil
}

func (s *InMemoryTokenStore) Delete(tokenID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, tokenID)
}
