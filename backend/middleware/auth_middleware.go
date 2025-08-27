package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDContextKey, int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
