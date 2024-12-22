package types

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type User struct {
	ID          string               `bson:"_id,omitempty"`
	Username    string               `bson:"username"`
	Email       string               `bson:"email"`
	Roles       []string             `bson:"roles"`
	ImageURL    string               `bson:"profileImageURL"`
	UsedStorage primitive.Decimal128 `bson:"usedStorage"`
}

type Claims struct {
	User User  `json:"user"`
	Exp  int64 `json:"exp"`
	jwt.RegisteredClaims
}
