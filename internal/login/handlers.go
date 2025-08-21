package login

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/rohitashvadangi/identity-server/internal/proto"
	"github.com/rohitashvadangi/identity-server/internal/stores"
)

type Login struct {
	authCodeStore stores.AuthCodeStore
	userStore     stores.UserStore
}

// ----- Users -----
var users = map[string]string{
	"alice": "password123",
	"bob":   "mypassword",
}

// /login - handles login and generates auth code
func (s *Login) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	scope := r.FormValue("scope")
	state := r.FormValue("state")

	if _, ok := s.userStore.ValidateCredentials(username, password); !ok {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate auth code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	ac := &proto.AuthCode{
		UserID:      username,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		Code:        code,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
	s.authCodeStore.Save(ac)

	// Redirect to client with code and state
	redirectURL := fmt.Sprintf("%s?code=%s", redirectURI, url.QueryEscape(code))
	if state != "" {
		redirectURL += "&state=" + url.QueryEscape(state)
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func NewLoginHandler(authCodeStore stores.AuthCodeStore, userStore stores.UserStore) *Login {
	return &Login{
		authCodeStore: authCodeStore,
		userStore:     userStore,
	}
}
