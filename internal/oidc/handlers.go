package oidc

import (
	"encoding/json"
	"net/http"
)

// Minimal UserInfo
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := map[string]interface{}{
		"sub":   "user123",
		"name":  "John Doe",
		"email": "john@example.com",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Minimal discovery
func DiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	base := "http://localhost:9090"
	resp := map[string]interface{}{
		"issuer":                                base,
		"authorization_endpoint":                base + "/authorize",
		"token_endpoint":                        base + "/token",
		"userinfo_endpoint":                     base + "/userinfo",
		"jwks_uri":                              base + "/.well-known/jwks.json",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Minimal JWKS
func JWKSHandler(w http.ResponseWriter, r *http.Request) {
	jwks := map[string]interface{}{
		"keys": []interface{}{}, // placeholder
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
