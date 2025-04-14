package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ulyanovikovak/PVZ_avito/internal/config"
)

type contextKey string

const (
	ContextKeyRole  contextKey = "role"
	ContextKeyToken contextKey = "token"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"message":"missing token"}`, http.StatusForbidden)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return config.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"message":"unauthorized"}`, http.StatusForbidden)
			return
		}

		role, _ := claims["role"].(string)
		ctx := context.WithValue(r.Context(), ContextKeyRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
