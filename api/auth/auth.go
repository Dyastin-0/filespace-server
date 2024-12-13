package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	user "filespace-backend/models"
)

func Handler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Bad request. Missing input: Email.", http.StatusBadRequest)
			return
		}

		if reqBody.Email == "" {
			http.Error(w, "Bad request. Missing input: Email.", http.StatusBadRequest)
			return
		}
		if reqBody.Password == "" {
			http.Error(w, "Bad request. Invalid input: Password.", http.StatusBadRequest)
			return
		}

		collection := client.Database("your_db_name").Collection("users")
		var user user.Model
		err := collection.FindOne(context.Background(), bson.M{"email": reqBody.Email}).Decode(&user)
		if err != nil {
			http.Error(w, "Account not found.", http.StatusNotFound)
			return
		}

		if !user.Verified {
			http.Error(w, "Verify your account.", http.StatusForbidden)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
		if err != nil {
			http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
			return
		}

		accessToken, err := generateToken(user, os.Getenv("ACCESS_TOKEN_KEY"), 15*time.Minute)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		newRefreshToken, err := generateToken(user, os.Getenv("REFRESH_TOKEN_KEY"), 24*time.Hour)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		cookies := r.Cookies()
		var newRefreshTokens []string
		if len(cookies) > 0 {
			for _, cookie := range cookies {
				if cookie.Name == "jwt" {
					foundToken := false
					for _, rt := range user.RefreshToken {
						if rt == cookie.Value {
							foundToken = true
							break
						}
					}
					if !foundToken {
						newRefreshTokens = []string{}
					}
					http.SetCookie(w, &http.Cookie{
						Name:     "jwt",
						Value:    "",
						HttpOnly: true,
						SameSite: http.SameSiteNoneMode,
						Secure:   true,
						MaxAge:   -1,
					})
				} else {
					newRefreshTokens = append(newRefreshTokens, cookie.Value)
				}
			}
		} else {
			newRefreshTokens = user.RefreshToken
		}

		newRefreshTokens = append(newRefreshTokens, newRefreshToken)
		_, err = collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"refreshToken": newRefreshTokens}})
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
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

		user.Password = ""
		user.RefreshToken = nil
		json.NewEncoder(w).Encode(map[string]interface{}{
			"accessToken": accessToken,
			"user":        user,
		})
	}
}

func generateToken(user user.Model, secret string, expiration time.Duration) (string, error) {
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
