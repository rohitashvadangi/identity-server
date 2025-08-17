package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// In-memory storage (demo)
var authCodes = map[string]string{} // code -> clientID

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	clientID := q.Get("client_id")
	redirectURI := q.Get("redirect_uri")
	scope := q.Get("scope")
	state := q.Get("state")

	if clientID == "" || redirectURI == "" {
		http.Error(w, "Missing client_id or redirect_uri", http.StatusBadRequest)
		return
	}

	// Show login form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<form method="POST" action="/login">
        <input type="hidden" name="client_id" value="%s">
        <input type="hidden" name="redirect_uri" value="%s">
        <input type="hidden" name="scope" value="%s">
        <input type="hidden" name="state" value="%s">
        Username: <input name="username"><br>
        Password: <input type="password" name="password"><br>
        <button type="submit">Login</button>
    </form>`, clientID, redirectURI, scope, state)
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
