package auth

import "github.com/dgrijalva/jwt-go"

type Body struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	AccessToken string   `json:"accessToken"`
	Email       string   `json:"user"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
}

type Claims struct {
	User struct {
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Roles    []string `json:"roles"`
		ID       string   `json:"_id"`
	} `json:"user"`
	Exp int64 `json:"exp"`
	jwt.StandardClaims
}
