package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
	mail "filespace/pkg/mail"
	mailTemplate "filespace/pkg/mail/template"
	hash "filespace/pkg/util/hash"
)

func Recover(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := types.RecoverBody{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request. Invalid format", http.StatusBadRequest)
			return
		}

		if reqBody.Token == "" || reqBody.NewPassword == "" {
			http.Error(w, "Bad request. Missing required fields", http.StatusBadRequest)
			return
		}

		claims := &types.Claims{}
		_, err = jwt.ParseWithClaims(reqBody.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("EMAIL_TOKEN_KEY")), nil
		})

		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		collection := client.Database("test").Collection("users")
		hashedPassword, _ := hash.Generate(reqBody.NewPassword)
		update := bson.M{"$set": bson.M{"password": hashedPassword}}
		filter := bson.M{"email": claims.User.Email}

		res, err := collection.UpdateOne(r.Context(), filter, update)

		if res == nil {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}

		if err != mongo.ErrNoDocuments && err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		options := mail.Message{
			To:      claims.User.Email,
			Subject: "Password recovery",
			Body: mailTemplate.Default(
				"Password recovery",
				"Your password has been successfully changed.",
				os.Getenv("BASE_CLIENT_URL")+"/sign-in",
				"Sign in now",
			),
		}

		mail.SendHTMLEmail(&options)

		w.WriteHeader(http.StatusOK)
	}
}
