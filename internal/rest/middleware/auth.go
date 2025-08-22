package middleware

import (
	"context"
	"cryptoserver/domain"
	"cryptoserver/internal/rest"
	"cryptoserver/pkg/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtConfig jwt.JWTConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				rest.WriteError(w, domain.ErrInvalidToken)
				return
			}

			tokenString, _ = strings.CutPrefix(tokenString, "Bearer ")
			claims, err := jwtConfig.ValidateToken(tokenString)
			if err != nil {
				rest.WriteError(w, domain.ErrInvalidToken)
				return
			}

			userID := (*claims)["sub"].(string)
			ctx := context.WithValue(r.Context(), "userID", userID)
			r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
