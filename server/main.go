package main

import (
	"log"
	"net/http"

	"github.com/rohitashvadangi/identity-server/internal/login"
	"github.com/rohitashvadangi/identity-server/internal/middleware"
	"github.com/rohitashvadangi/identity-server/internal/oauth"
	"github.com/rohitashvadangi/identity-server/internal/oidc"
	"github.com/rohitashvadangi/identity-server/internal/stores"
	"github.com/rohitashvadangi/identity-server/internal/utils"
)

func main() {
	mux := http.NewServeMux()
	userStore := stores.NewUserStore(utils.CreateUsers())

	authCodeStore := stores.NewAuthCodeStore()
	tokenStore := stores.NewTokenStore()
	loginHandler := login.NewLoginHandler(authCodeStore, userStore)
	oauthHandler := oauth.NewOauthHandler(authCodeStore, tokenStore)

	userHandler := oidc.NewUserHandler(userStore)
	oidc.InitJWKS("../keys/public.pem")
	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// OAuth endpoints
	mux.HandleFunc("/authorize", oauthHandler.AuthorizeHandler)
	mux.HandleFunc("/token", oauthHandler.TokenHandler)
	mux.HandleFunc("/introspect", oauthHandler.IntrospectHandler)
	mux.HandleFunc("/revoke", oauthHandler.RevokeHandler)

	// OIDC endpoints
	var validator stores.TokenValidator = tokenStore

	mux.Handle("/userinfo",
		middleware.AuthMiddleware(validator)(
			http.HandlerFunc(userHandler.UserInfoHandler),
		),
	)

	mux.HandleFunc("/.well-known/openid-configuration", oidc.DiscoveryHandler)
	mux.HandleFunc("/.well-known/jwks.json", oidc.JWKSHandler)

	//login demo
	mux.HandleFunc("/login", loginHandler.LoginHandler)

	log.Println("Identity Server running on :9090")
	srv := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	log.Println("Identity Server running on :9090")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
