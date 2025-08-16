package oauth

import (
	"encoding/json"
	"net/http"
)

// In-memory storage (demo)
var authCodes = map[string]string{} // code -> clientID

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	// Very minimal: accept ?client_id=&redirect_uri=
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	if clientID == "" || redirectURI == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	code := "authcode123" // demo code
	authCodes[code] = clientID

	// Redirect with code
	http.Redirect(w, r, redirectURI+"?code="+code, http.StatusFound)
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	grantType := r.FormValue("grant_type")
	if grantType != "authorization_code" {
		http.Error(w, "unsupported grant_type", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	clientID := authCodes[code]
	if clientID == "" {
		http.Error(w, "invalid code", http.StatusBadRequest)
		return
	}

	// Issue dummy access & ID token
	resp := map[string]interface{}{
		"access_token":  "access123",
		"token_type":    "Bearer",
		"expires_in":    3600,
		"id_token":      "idtoken123",
		"refresh_token": "refresh123",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
