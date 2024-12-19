package auth

import (
	"encoding/json"
	"log"

	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	types "filespace/internal/auth/types"
	user "filespace/internal/models/user"
	utils "filespace/utils"
)

func Handler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody = types.Body{}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Fatal(err)
			http.Error(w, "Bad request.", http.StatusBadRequest)
			return
		}

		if reqBody.Email == "" || reqBody.Password == "" {
			http.Error(w, "Bad request. Missing required fields.", http.StatusBadRequest)
			return
		}

		collection := client.Database("test").Collection("users")
		var user user.Model
		err := collection.FindOne(r.Context(), bson.M{"email": reqBody.Email}).Decode(&user)
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

		accessToken, err := utils.GenerateToken(user, os.Getenv("ACCESS_TOKEN_KEY"), 15*time.Minute)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		newRefreshToken, err := utils.GenerateToken(user, os.Getenv("REFRESH_TOKEN_KEY"), 24*time.Hour)
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
		_, err = collection.UpdateOne(r.Context(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"refreshToken": newRefreshTokens}})
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

		response := types.Response{
			AccessToken: accessToken,
			Email:       user.Email,
			Username:    user.Username,
			Roles:       user.Roles,
		}
		json.NewEncoder(w).Encode(response)
	}
}
