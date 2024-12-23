package file

import (
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "filespace/internal/auth/type"
)

func Post(storageClient *storage.Client, mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		id := claims.User.ID
		bucket := storageClient.Bucket(bucketName)

		r.ParseMultipartForm(10 << 20)
		files := r.MultipartForm.File["files"]
		path := r.FormValue("path")
		folder := r.FormValue("folder")
		size := r.Context().Value("size")

		if files != nil {
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					http.Error(w, "Failed to open file", http.StatusInternalServerError)
					return
				}

				defer file.Close()

				fileName := fmt.Sprintf("%s/%s", id, fileHeader.Filename)
				if path != "" {
					fileName = fmt.Sprintf("%s/%s/%s", id, path, fileHeader.Filename)
				}

				newFile := bucket.Object(fileName)
				writer := newFile.NewWriter(r.Context())
				writer.ContentType = fileHeader.Header.Get("Content-Type")

				fileBytes, err := io.ReadAll(file)
				if err != nil {
					http.Error(w, "Error reading file", http.StatusInternalServerError)
					return
				}

				if _, err := writer.Write(fileBytes); err != nil {
					http.Error(w, "Error writing file to GCS", http.StatusInternalServerError)
					return
				}

				if err := writer.Close(); err != nil {
					http.Error(w, "Error closing GCS writer", http.StatusInternalServerError)
					return
				}

				attrsToUpdate := storage.ObjectAttrsToUpdate{
					Metadata: map[string]string{
						"owner": id,
					},
				}
				if _, err := newFile.Update(r.Context(), attrsToUpdate); err != nil {
					http.Error(w, "Error setting metadata", http.StatusInternalServerError)
					return
				}
			}

			collection := mongoClient.Database("test").Collection("users")
			filter := bson.M{"email": claims.User.Email}

			update := bson.M{"$inc": bson.M{"usedStorage": size}}
			if _, err := collection.UpdateOne(r.Context(), filter, update); err != nil {
				http.Error(w, "Error updating user storage", http.StatusInternalServerError)
				return
			}

		} else if folder != "" {
			folderName := fmt.Sprintf("%s/%s/", id, folder)
			if path != "" {
				folderName = fmt.Sprintf("%s/%s/%s/", id, path, folder)
			}
			newFolder := bucket.Object(folderName)
			if err := newFolder.NewWriter(r.Context()).Close(); err != nil {
				http.Error(w, "Error creating folder", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
