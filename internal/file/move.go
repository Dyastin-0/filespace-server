package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"cloud.google.com/go/storage"

	types "filespace/internal/auth/type"
	fileTypes "filespace/internal/file/type"
)

func Move(client *storage.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := fileTypes.MoveBody{}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Bad request. Invalid format", http.StatusBadRequest)
			return
		}

		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		prefix := claims.User.ID

		srcFolderPath := fmt.Sprintf("%s/%s", prefix, reqBody.File.Path)
		newFolderPath := fmt.Sprintf("%s/%s/%s", prefix, reqBody.TargetPath, path.Base(reqBody.File.Path))
		if reqBody.TargetPath == "" {
			newFolderPath = fmt.Sprintf("%s/%s", prefix, path.Base(reqBody.File.Path))
		}

		bucket := client.Bucket(bucketName)
		query := &storage.Query{Prefix: srcFolderPath}
		objIter := bucket.Objects(r.Context(), query)

		for {
			attrs, err := objIter.Next()
			if err != nil {
				break
			}

			relativePath := attrs.Name[len(srcFolderPath):]
			newFileName := newFolderPath + relativePath

			newObj := bucket.Object(newFileName)

			if _, err := newObj.CopierFrom(bucket.Object(attrs.Name)).Run(r.Context()); err != nil {
				http.Error(w, "Error copying file", http.StatusInternalServerError)
				return
			}

			if err := bucket.Object(attrs.Name).Delete(r.Context()); err != nil {
				http.Error(w, "Error deleting original file", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Folder and its contents moved successfully."))
	}
}
