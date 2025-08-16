package middleware

import (
	"context"
	"my-notes-app/pkg/jwt"
	"net/http"
	"strings"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.ParseToken(tokenStr, secret)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
