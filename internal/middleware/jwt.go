package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	types "filespace/internal/auth/type"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		claims := &types.Claims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("EMAIL_TOKEN_KEY")), nil
		})

		if err != nil {
			http.Error(w, "Forbidden.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})

}
