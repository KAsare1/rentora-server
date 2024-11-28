package middleware

import (
	"context"
	"errors"
	"net/http"
	"rentora-go/internal/service"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserContextKey ContextKey = "user"

// AuthMiddleware validates the JWT token and attaches user info to the request context.
func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate token
			claims := &service.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Attach claims to request context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user info from the request context.
func GetUserFromContext(ctx context.Context) (*service.Claims, error) {
	claims, ok := ctx.Value(UserContextKey).(*service.Claims)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return claims, nil
}
