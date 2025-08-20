package middleware

import (
	"context"
	"cryptoserver/domain"
	"cryptoserver/internal/rest"
	"cryptoserver/pkg/jwt"
	"net/http"
)

func AuthMiddleware(jwtConfig jwt.JWTConfig, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			rest.WriteError(w, domain.ErrInvalidToken)
			return
		}
		claims, err := jwtConfig.ValidateToken(tokenString)
		if err != nil {
			rest.WriteError(w, domain.ErrInvalidToken)
			return
		}
		userID := (*claims)["sub"].(string)
		ctx := context.WithValue(r.Context(), "userID", userID)
		r.WithContext(ctx)
		next(w, r)
	}

}
