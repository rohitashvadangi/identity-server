package main

import (
	"log"
	"net/http"

	"github.com/rohitashvadangi/identity-server/internal/login"
	"github.com/rohitashvadangi/identity-server/internal/oauth"
	"github.com/rohitashvadangi/identity-server/internal/oidc"
	"github.com/rohitashvadangi/identity-server/internal/stores"
)

func main() {
	mux := http.NewServeMux()
	loginHander := login.NewLoginHandler(stores.NewAuthCodeStore())

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// OAuth endpoints
	mux.HandleFunc("/authorize", oauth.AuthorizeHandler)
	mux.HandleFunc("/token", oauth.TokenHandler)

	// OIDC endpoints
	mux.HandleFunc("/userinfo", oidc.UserInfoHandler)
	mux.HandleFunc("/.well-known/openid-configuration", oidc.DiscoveryHandler)
	mux.HandleFunc("/.well-known/jwks.json", oidc.JWKSHandler)
		mux.HandleFunc("/login", loginHander.LoginHandler)


	log.Println("Identity Server running on :9090")
	log.Fatal(http.ListenAndServe(":9090", mux))
}
