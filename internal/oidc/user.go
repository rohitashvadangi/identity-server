package oidc

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rohitashvadangi/identity-server/internal/stores"
)

type UserHandler struct {
	userStore stores.UserStore
}

// Minimal UserInfo
func (u *UserHandler) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		json.NewEncoder(w).Encode(errors.New("Not found"))
	}
	result, err := u.userStore.FindByUsername(userID)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	user := map[string]interface{}{
		"sub":        result.UserName,
		"familyName": result.FamilyName,
		"givenName":  result.GivenName,
		"email":      result.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func NewUserHandler(userStore stores.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}
