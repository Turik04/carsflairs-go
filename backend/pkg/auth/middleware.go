package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Проверяем роль
		if claims.Role != requiredRole {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		// Если все ок — вызываем следующий обработчик
		next(w, r)
	}
}
