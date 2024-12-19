package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	auth "filespace/internal/auth/types"
	user "filespace/internal/models/user"
	mail "filespace/pkg/mail"
	mailTemplate "filespace/pkg/mail/templates"
	mailTypes "filespace/pkg/mail/types"
	utils "filespace/utils"
)

func Signup(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody = auth.SignupBody{}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Bad request. Invalid format.", http.StatusBadRequest)
			return
		}

		if reqBody.Email == "" || reqBody.Password == "" || reqBody.Username == "" {
			http.Error(w, "Bad request. Missing required fields.", http.StatusBadRequest)
			return
		}

		collection := client.Database("test").Collection("users")
		err := collection.FindOne(context.Background(), bson.M{"email": reqBody.Email}).Err()

		if err == nil {
			http.Error(w, "Email already used.", http.StatusConflict)
			return
		} else if err != mongo.ErrNoDocuments {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		err = collection.FindOne(context.Background(), bson.M{"username": reqBody.Username}).Err()
		if err == nil {
			http.Error(w, "Username already used.", http.StatusConflict)
			return
		} else if err != mongo.ErrNoDocuments {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		hashedPassword, herr := utils.GenerateHash(reqBody.Password)
		if herr != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		user := user.Model{
			Email:    reqBody.Email,
			Password: hashedPassword,
			Username: reqBody.Username,
			Verified: false,
			Roles:    []string{"122602"},
		}

		verificationToken, err := utils.GenerateToken(user, os.Getenv("EMAIL_TOKEN_KEY"), 5*time.Minute)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		user.VerificationToken = verificationToken

		_, err = collection.InsertOne(context.Background(), user)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		message := &mailTypes.Message{
			To:          user.Email,
			Subject:     "Verify your account",
			ContentType: mailTypes.HTMLTextEmail,
			Body: mailTemplate.Default(
				"Verification",
				"Hi, "+user.Username+". Thank you for signing up. Please verify your account by clicking the link below. if you did not sign up, please ignore this email.",
				os.Getenv("BASE_CLIENT_URL")+"/account/verification?t="+verificationToken,
				"Verify Account"),
		}

		err = mail.SendHTMLEmail(message)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
