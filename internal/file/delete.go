package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/bson"
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
		totalDeletedSize := int64(0)

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
					totalDeletedSize += objAttrs.Size
					if err := obj.Delete(r.Context()); err != nil {
						http.Error(w, "Failed to delete file in folder", http.StatusInternalServerError)
						return
					}
				}
			} else {
				objAttrs, err := obj.Attrs(r.Context())
				if err != nil {
					if err == storage.ErrObjectNotExist {
						continue
					}
					http.Error(w, "Failed to get file attributes", http.StatusInternalServerError)
					return
				}
				totalDeletedSize += objAttrs.Size

				if err := obj.Delete(r.Context()); err != nil && err != storage.ErrObjectNotExist {
					http.Error(w, "Failed to delete file", http.StatusInternalServerError)
					return
				}
			}
		}

		collection := mongoClient.Database("test").Collection("users")
		filter := bson.M{"email": claims.User.Email}
		update := bson.M{"$inc": bson.M{"usedStorage": -totalDeletedSize}}
		if _, err := collection.UpdateOne(r.Context(), filter, update); err != nil {
			http.Error(w, "Error updating user storage", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
