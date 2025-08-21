package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rohitashvadangi/identity-server/internal/stores"
)

func AuthMiddleware(tokenStore stores.TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			tokenDetails, ok := tokenStore.ValidateOpaqueToken(token)
			if !ok {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", tokenDetails.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
