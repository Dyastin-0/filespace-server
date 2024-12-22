package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
	fileTypes "filespace/internal/file/type"
)

func Delete(storageClient *storage.Client, mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := fileTypes.DeleteBody{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id := claims.User.ID

		if len(reqBody.Files) == 0 {
			http.Error(w, "File names are required", http.StatusBadRequest)
			return
		}

		bucket := storageClient.Bucket(bucketName)

		for _, file := range reqBody.Files {
			filePath := fmt.Sprintf("%s/%s", id, file)
			obj := bucket.Object(filePath)

			if strings.HasSuffix(file, "/") {
				query := &storage.Query{Prefix: filePath}
				it := bucket.Objects(r.Context(), query)
				for {
					objAttrs, err := it.Next()
					if err != nil {
						break
					}

					obj := bucket.Object(objAttrs.Name)
					if err := obj.Delete(r.Context()); err != nil {
						http.Error(w, "Failed to delete file in folder", http.StatusInternalServerError)
						return
					}
				}
			} else {
				if err := obj.Delete(r.Context()); err != nil {
					http.Error(w, "Failed to delete file", http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
