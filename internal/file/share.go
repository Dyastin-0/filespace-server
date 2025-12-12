package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
	fileTypes "filespace/internal/file/type"
	mail "filespace/pkg/mail"
	mailTemplate "filespace/pkg/mail/template"
)

func Share(storageClient *storage.Client, mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := fileTypes.ShareBody{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request. Invalid format", http.StatusBadRequest)
			return
		}

		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id := claims.User.ID

		prefix := fmt.Sprintf("%s/%s", id, reqBody.File)
		bucket := storageClient.Bucket(bucketName)

		obj := bucket.Object(prefix)

		attrs, err := obj.Attrs(r.Context())

		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		signedURL, err := bucket.SignedURL(attrs.Name, &storage.SignedURLOptions{
			Method:  "GET",
			Expires: time.Now().Add(time.Duration(reqBody.Exp.Value) * time.Millisecond),
		})

		if err != nil {
			http.Error(w, "Failed to generate signed URL", http.StatusInternalServerError)
			return
		}

		user := types.User{}
		collection := mongoClient.Database("test").Collection("users")
		err = collection.FindOne(r.Context(), bson.M{"email": reqBody.Email}).Decode(&user)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		options := mail.Message{
			To:          reqBody.Email,
			Subject:     "File sharing",
			ContentType: mail.HTMLTextEmail,
			Body: mailTemplate.Default(
				user.Username+" shared a file with you",
				"You can access the file using the link below. The link will expire in "+reqBody.Exp.Str+".",
				signedURL,
				"Access file",
			),
		}

		mail.SendHTMLEmail(&options)

		w.WriteHeader(http.StatusOK)
	}
}
