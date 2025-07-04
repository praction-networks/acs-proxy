package middleware

import (
	"net/http"
	"strings"

	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/helpers"
)

func APIKeyAuthMiddleware(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				helpers.HandleAppError(w, appError.New(appError.UnauthorizedAccess, "Missing or invalid Authorization header", http.StatusUnauthorized, nil))
				return
			}

			apiKey := strings.TrimPrefix(authHeader, "Bearer ")
			if apiKey != expectedKey {
				helpers.HandleAppError(w, appError.New(appError.UnauthorizedAccess, "Invalid API key", http.StatusUnauthorized, nil))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
