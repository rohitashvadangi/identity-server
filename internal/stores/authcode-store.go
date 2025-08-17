package stores

import (
	"errors"
	"sync"
	"time"

	"github.com/rohitashvadangi/identity-server/internal/proto"
)

type AuthCodeStore interface {
	Save(code *proto.AuthCode)

	Get(code string) (*proto.AuthCode, error)

	Delete(code string)
}

type AuthCodeStoreImpl struct {
	codes map[string]*proto.AuthCode
	mu    sync.RWMutex
}

func NewAuthCodeStore() *AuthCodeStoreImpl {
	return &AuthCodeStoreImpl{
		codes: make(map[string]*proto.AuthCode),
	}
}

func (s *AuthCodeStoreImpl) Save(code *proto.AuthCode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[code.Code] = code
}

func (s *AuthCodeStoreImpl) Get(code string) (*proto.AuthCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.codes[code]
	if !ok || time.Now().After(c.ExpiresAt) {
		return nil, errors.New("auth code not found or expired")
	}
	return c, nil
}

func (s *AuthCodeStoreImpl) Delete(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.codes, code)
}
