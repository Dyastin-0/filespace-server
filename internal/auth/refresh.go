package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	authTypes "filespace/internal/auth/type"
	user "filespace/internal/model/user"
	token "filespace/pkg/util/token"
)

func Refresh(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("rt")

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		refreshToken := cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:     "rt",
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			MaxAge:   -1,
			Path:     "/",
		})

		collection := client.Database("test").Collection("users")
		var user user.Model
		err = collection.FindOne(context.Background(), bson.M{"refreshToken": refreshToken}).Decode(&user)
		if err != nil {
			claims := &authTypes.Claims{}
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("REFRESH_TOKEN_KEY")), nil
			})
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			_, err = collection.UpdateOne(context.Background(), bson.M{"email": claims.Subject}, bson.M{"$set": bson.M{"refreshToken": []string{}}})
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		newRefreshTokens := []string{}
		for _, rt := range user.RefreshToken {
			if rt != refreshToken {
				newRefreshTokens = append(newRefreshTokens, rt)
			}
		}

		claims := &authTypes.Claims{}
		_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REFRESH_TOKEN_KEY")), nil
		})
		if err != nil {
			_, err = collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"refreshToken": newRefreshTokens}})
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if user.Email != claims.User.Email {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		accessToken, err := token.Generate(&user, os.Getenv("ACCESS_TOKEN_KEY"), 15*time.Minute)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		newRefreshToken, err := token.Generate(&user, os.Getenv("REFRESH_TOKEN_KEY"), 24*time.Hour)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		newRefreshTokens = append(newRefreshTokens, newRefreshToken)
		_, err = collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"refreshToken": newRefreshTokens}})
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "rt",
			Value:    newRefreshToken,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			MaxAge:   24 * 60 * 60,
			Path:     "/",
		})

		response := authTypes.Response{
			AccessToken: accessToken,
			User: authTypes.User{
				Username:    user.Username,
				Email:       user.Email,
				Roles:       user.Roles,
				ImageURL:    user.ImageURL,
				UsedStorage: user.UsedStorage,
			},
		}

		json.NewEncoder(w).Encode(response)
	}
}
