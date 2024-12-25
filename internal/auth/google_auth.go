package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	user "filespace/internal/model/user"
	token "filespace/pkg/util/token"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func Google() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := googleOAuthConfig.AuthCodeURL("state")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func GoogleCallback(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		tkn, err := googleOAuthConfig.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			fmt.Printf("Error exchanging token: %v\n", err)
			return
		}

		httpClient := googleOAuthConfig.Client(r.Context(), tkn)
		userInfoResponse, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			fmt.Printf("Error fetching user info: %v\n", err)
			return
		}
		defer userInfoResponse.Body.Close()

		userInfoBytes, err := io.ReadAll(userInfoResponse.Body)
		if err != nil {
			http.Error(w, "Failed to read user info", http.StatusInternalServerError)
			return
		}

		var userInfo map[string]interface{}
		if err := json.Unmarshal(userInfoBytes, &userInfo); err != nil {
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}

		collection := client.Database("test").Collection("users")
		filter := bson.M{"googleId": userInfo["id"].(string)}
		existingUser := user.Model{}

		err = collection.FindOne(r.Context(), filter).Decode(&existingUser)
		if err == mongo.ErrNoDocuments {
			newUser := user.Model{
				Username: userInfo["name"].(string),
				Email:    userInfo["email"].(string),
				GoogleID: userInfo["id"].(string),
				Verified: true,
				Roles:    []string{"122602"},
				Created:  time.Now(),
			}

			_, err := collection.InsertOne(r.Context(), newUser)
			if err != nil {
				http.Error(w, "Failed to insert new user", http.StatusInternalServerError)
				fmt.Printf("Error inserting new user: %v\n", err)
				return
			}
			existingUser = newUser
		} else if err != nil {
			http.Error(w, "Failed to find user", http.StatusInternalServerError)
			fmt.Printf("Error finding user: %v\n", err)
			return
		}

		refreshToken, err := token.Generate(&existingUser, os.Getenv("REFRESH_TOKEN_KEY"), 24*time.Hour)
		if err != nil {
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			fmt.Printf("Error generating refresh token: %v\n", err)
			return
		}

		update := bson.M{"$push": bson.M{"refreshToken": refreshToken}}
		_, err = collection.UpdateOne(r.Context(), filter, update)
		if err != nil {
			http.Error(w, "Failed to update user with refresh token", http.StatusInternalServerError)
			fmt.Printf("Error updating user: %v\n", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    refreshToken,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			MaxAge:   24 * 60 * 60,
		})

		redirectURL := os.Getenv("BASE_CLIENT_URL") + "/home"
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}
