package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	auth "filespace/internal/auth/types"

	user "filespace/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

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
		if err == nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		user.VerificationToken = verificationToken

		_, err = collection.InsertOne(context.Background(), user)
		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
