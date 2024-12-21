package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
	user "filespace/internal/model/user"
	mail "filespace/pkg/mail"
	mailTemplate "filespace/pkg/mail/template"
	token "filespace/pkg/util/token"
)

func SendRecovery(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := types.SendRecoveryBody{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request. Invalid format.", http.StatusBadRequest)
			return
		}

		if reqBody.Email == "" {
			http.Error(w, "Bad request. Missing required fields.", http.StatusBadRequest)
			return
		}

		collection := client.Database("test").Collection("users")
		var user user.Model
		err = collection.FindOne(r.Context(), bson.M{"email": reqBody.Email}).Decode(&user)

		if err != nil {
			http.Error(w, "Account not found.", http.StatusNotFound)
			return
		}

		token, err := token.Generate(&user, os.ExpandEnv("EMAIL_TOKEN_KEY"), time.Hour*24)

		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		options := mail.Message{
			To:          reqBody.Email,
			Subject:     "Password recovery",
			ContentType: mail.HTMLTextEmail,
			Body: mailTemplate.Default(
				"Password Recovery",
				"Click the link below to recover your password",
				os.Getenv("BASE_CLIENT_URL")+"/auth/recover?t="+token,
				"Recover password",
			),
		}

		err = mail.SendHTMLEmail(&options)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
