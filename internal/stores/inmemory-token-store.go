package stores

import (
	"errors"
	"sync"
	"time"

	"github.com/rohitashvadangi/identity-server/internal/proto"
)

type InMemoryTokenStore struct {
	tokens        map[string]*proto.Token
	refreshTokens map[string]*proto.RefreshToken
	mu            sync.RWMutex
}

// Ensure InMemoryTokenStore implements TokenStore
var _ TokenStore = (*InMemoryTokenStore)(nil)

func NewTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		tokens: make(map[string]*proto.Token), refreshTokens: make(map[string]*proto.RefreshToken),
	}
}

func (s *InMemoryTokenStore) Save(token *proto.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token.TokenID] = token
}

func (s *InMemoryTokenStore) SaveRefresh(refreshToken *proto.RefreshToken) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshTokens[refreshToken.ID] = refreshToken
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

func (s *InMemoryTokenStore) GetByRefresh(refresh string) (*proto.RefreshToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.refreshTokens[refresh]
	if !ok || time.Now().After(t.ExpiresAt) {
		return nil, errors.New("refresh token not found or expired")
	}
	return t, nil
}

// GenerateJWT signs a token with RSA private key
func (s *InMemoryTokenStore) ValidateOpaqueToken(token string) (*proto.Token, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tokens[token]
	if !ok {
		return nil, false
	}

	if time.Now().After(t.ExpiresAt) {
		delete(s.tokens, token)

		return nil, false
	}
	return t, true
}
