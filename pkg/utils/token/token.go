package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	authTypes "filespace/internal/auth/types"
	user "filespace/internal/models/user"
)

func Generate(user user.Model, secret string, expiration time.Duration) (string, error) {
	claims := authTypes.Claims{
		User: struct {
			Username string   `json:"username"`
			Email    string   `json:"email"`
			Roles    []string `json:"roles"`
			ID       string   `json:"_id"`
		}{
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.Roles,
			ID:       user.ID,
		},
		Exp: time.Now().Add(expiration).Unix(),
		StandardClaims: jwt.StandardClaims{
			Issuer:   "Filespace",
			Subject:  user.ID,
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
