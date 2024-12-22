package middleware

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
)

const MAX_STORAGE_LIMIT = 1 * 1024 * 1024 * 1024 // 1 GB

type User struct {
	ID          string `bson:"_id"`
	UsedStorage int64  `bson:"usedStorage"`
}

func CheckStorageCapacity(client *mongo.Client) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*types.Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			id := claims.User.ID

			r.ParseMultipartForm(10 << 20)
			files := r.MultipartForm.File["files"]

			if len(files) == 0 {
				http.Error(w, "No files to upload.", http.StatusBadRequest)
				return
			}

			collection := client.Database("test").Collection("users")
			var user User
			err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
			if err != nil && err != mongo.ErrNoDocuments {
				http.Error(w, "Error checking storage limit.", http.StatusInternalServerError)
				return
			}

			usedStorage := user.UsedStorage
			totalUploadSize := int64(0)
			for _, fileHeader := range files {
				totalUploadSize += fileHeader.Size
			}

			if usedStorage+totalUploadSize > MAX_STORAGE_LIMIT {
				http.Error(w, "Storage limit exceeded. Cannot upload files.", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "size", totalUploadSize)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
