package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	types "filespace/internal/auth/type"
	usr "filespace/internal/model/user"
	token "filespace/pkg/util/token"
)

func Verify(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		var reqQuery = types.VerifyQuery{}

		reqQuery.Token = query.Get("t")

		if reqQuery.Token == "" {
			http.Error(w, "Bad request. Missing required fields.", http.StatusBadRequest)
			return
		}

		claims := &types.Claims{}
		_, err := jwt.ParseWithClaims(reqQuery.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("EMAIL_TOKEN_KEY")), nil
		})

		if err != nil {
			http.Error(w, "Invalid or expired token.", http.StatusUnauthorized)
			return
		}

		collection := client.Database("test").Collection("users")
		update := bson.M{"$set": bson.M{"verified": true, "verificationToken": ""}}
		filter := bson.M{"email": claims.User.Email}
		ops := options.Update().SetUpsert(true)
		res, err := collection.UpdateOne(r.Context(), filter, update, ops)

		if res == nil {
			http.Error(w, "Account not found.", http.StatusNotFound)
			return
		}

		if err != mongo.ErrNoDocuments && err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		user := usr.Model{}
		err = collection.FindOne(r.Context(), filter).Decode(&user)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		accessToken, err := token.Generate(&user, os.Getenv("ACCESS_TOKEN_KEY"), 15*time.Minute)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		refreshToken, err := token.Generate(&user, os.Getenv("REFRESH_TOKEN_KEY"), 24*time.Hour)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		filter = bson.M{"email": user.Email}
		update = bson.M{"$set": bson.M{"refreshToken": refreshToken}}
		_, err = collection.UpdateOne(r.Context(), filter, update)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
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

		response := types.Response{
			AccessToken: accessToken,
			Email:       user.Email,
			Username:    user.Username,
			Roles:       user.Roles,
		}

		json.NewEncoder(w).Encode(response)
	}
}
