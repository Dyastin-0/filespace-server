package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	authTypes "filespace/internal/auth/type"
	user "filespace/internal/model/user"
)

func Generate(user *user.Model, secret string, expiration time.Duration) (string, error) {
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
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Filespace",
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
