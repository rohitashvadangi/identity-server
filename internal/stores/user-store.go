package stores

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rohitashvadangi/identity-server/internal/proto"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	// Find a user by username
	FindByUsername(username string) (*proto.User, error)

	//Find by ID
	FindById(Id string) (*proto.User, error)

	// Save a new user
	Save(user *proto.User) error

	ValidateCredentials(username, password string) (*proto.User, bool)
}

type UserStoreImpl struct {
	usr map[string]*proto.User
	mu  sync.RWMutex
}

func (s *UserStoreImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserStoreImpl) ValidateCredentials(username, password string) (*proto.User, bool) {
	user, err := s.FindByUsername(username)
	if err != nil {
		return nil, false
	}

	if user.UserName == username && s.CheckPasswordHash(password, user.UserPwdHash) {
		return user, true
	}

	return nil, false
}

func (s *UserStoreImpl) Save(user *proto.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.usr[user.Id] = user
	return nil
}

func (s *UserStoreImpl) FindByUsername(username string) (*proto.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, value := range s.usr {
		if value.UserName == username {
			return value, nil
		}
	}

	return nil, errors.New("username not found ")

}

func (s *UserStoreImpl) FindById(Id string) (*proto.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.usr[Id]
	if !ok {
		fmt.Println("User not found for Id ", Id)
		return nil, errors.New("userId not found ")
	}
	return c, nil
}

func NewUserStore(users map[string]*proto.User) *UserStoreImpl {
	return &UserStoreImpl{
		usr: users,
	}
}
