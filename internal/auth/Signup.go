package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	auth "filespace/internal/auth/type"
	user "filespace/internal/model/user"
	mail "filespace/pkg/mail"
	mailTemplate "filespace/pkg/mail/template"
	hash "filespace/pkg/util/hash"
	token "filespace/pkg/util/token"
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

		hashedPassword, herr := hash.Generate(reqBody.Password)
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

		verificationToken, err := token.Generate(&user, os.Getenv("EMAIL_TOKEN_KEY"), 15*time.Minute)
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

		options := mail.Message{
			To:          user.Email,
			Subject:     "Verify your account",
			ContentType: mail.HTMLTextEmail,
			Body: mailTemplate.Default(
				"Verification",
				"Hi, "+user.Username+". Thank you for signing up. Please verify your account by clicking the link below. if you did not sign up, please ignore this email.",
				os.Getenv("BASE_CLIENT_URL")+"/auth/verify?t="+verificationToken,
				"Verify Account"),
		}

		err = mail.SendHTMLEmail(&options)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
