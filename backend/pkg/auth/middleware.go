package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

// Middleware — просто проверка токена
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserID, claims.ID)
		ctx = context.WithValue(ctx, ContextRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware — проверка роли, и только потом handler
func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(ContextRole).(string)
		if !ok || role != requiredRole {
			http.Error(w, "Сенде доступ жоккой інім", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}
