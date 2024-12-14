package token

import (
	user "filespace-backend/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Generate(user user.Model, secret string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user": map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"roles":    user.Roles,
			"_id":      user.ID,
		},
		"exp": time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
