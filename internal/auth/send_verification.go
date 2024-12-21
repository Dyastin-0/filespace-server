package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
	usr "filespace/internal/model/user"
	mail "filespace/pkg/mail"
	template "filespace/pkg/mail/template"
	token "filespace/pkg/util/token"
)

func SendVerification(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody = types.VerificationBody{}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Bad request. Invalid format.", http.StatusBadRequest)
			return
		}

		if reqBody.Email == "" {
			http.Error(w, "Bad request. Missing required fields.", http.StatusBadRequest)
			return
		}

		collection := client.Database("test").Collection("users")
		var user usr.Model
		err := collection.FindOne(r.Context(), bson.M{"email": reqBody.Email}).Decode(&user)

		if err != nil {
			http.Error(w, "Account not found.", http.StatusNotFound)
			return
		}

		verificationToken, err := token.Generate(&user, os.Getenv("EMAIL_TOKEN_KEY"), 15*time.Minute)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		options := mail.Message{
			To:          reqBody.Email,
			Subject:     "Email Verification",
			ContentType: mail.HTMLTextEmail,
			Body: template.Default(
				"Email Verification",
				"Click the link below to verify your email.",
				os.Getenv("BASE_CLIENT_URL")+"/auth/verify?t="+verificationToken,
				"Verify Email"),
		}

		err = mail.SendHTMLEmail(&options)

		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
