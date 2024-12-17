package refresh

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	user "filespace/models"
	authTypes "filespace/types/auth"
	"filespace/types/refresh"
	token "filespace/utils"
)

func Handler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		refreshToken := cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			MaxAge:   -1,
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
				http.Error(w, "Forbidden!", http.StatusForbidden)
				return
			}
			_, err = collection.UpdateOne(context.Background(), bson.M{"email": claims.Subject}, bson.M{"$set": bson.M{"refreshToken": []string{}}})
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			http.Error(w, "Forbidden!!", http.StatusForbidden)
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
			http.Error(w, "Forbidden!!!", http.StatusForbidden)
			return
		}

		if user.Email != claims.User.Email {
			http.Error(w, "Forbidden!!!!", http.StatusForbidden)
			return
		}

		accessToken, err := token.Generate(user, os.Getenv("ACCESS_TOKEN_SECRET"), 15*time.Minute)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		newRefreshToken, err := token.Generate(user, os.Getenv("REFRESH_TOKEN_KEY"), 24*time.Hour)
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
			Name:     "jwt",
			Value:    newRefreshToken,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			MaxAge:   24 * 60 * 60,
		})

		response := refresh.Response{
			AccessToken: accessToken,
		}

		json.NewEncoder(w).Encode(response)
	}
}
