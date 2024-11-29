// internal/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"rentora-go/internal/service"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get the Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
                return
            }

            // Bearer token validation
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
                return
            }

            tokenString := parts[1]

            // Parse and validate the token
            token, err := jwt.ParseWithClaims(tokenString, &service.Claims{}, func(token *jwt.Token) (interface{}, error) {
                return jwtKey, nil
            })

            if err != nil {
                http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
                return
            }

            // Extract claims
            claims, ok := token.Claims.(*service.Claims)
            if !ok || !token.Valid {
                http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
                return
            }

            // Add user ID to context
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}