package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	types "filespace/internal/auth/type"
	fileTypes "filespace/internal/file/type"

	"cloud.google.com/go/storage"
)

const bucketName = "filespace-bucket"

func Get(client *storage.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id := claims.User.ID

		prefix := fmt.Sprintf("%s/", id)
		bucket := client.Bucket(bucketName)

		query := &storage.Query{Prefix: prefix}
		it := bucket.Objects(r.Context(), query)

		var filesMetaData []fileTypes.Metadata

		for {
			attrs, err := it.Next()
			if err == storage.ErrBucketNotExist {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if err != nil {
				break
			}

			signedURL, err := bucket.SignedURL(attrs.Name, &storage.SignedURLOptions{
				Method:  "GET",
				Expires: time.Now().Add(15 * time.Minute),
			})
			if err != nil {
				http.Error(w, "Failed to generate signed URL", http.StatusInternalServerError)
				return
			}

			filesMetaData = append(filesMetaData, fileTypes.Metadata{
				Owner:       attrs.Metadata["owner"],
				Name:        attrs.Name[len(prefix):],
				Link:        signedURL,
				Size:        attrs.Size,
				Type:        getFileType(attrs.ContentType),
				Updated:     attrs.Updated,
				ContentType: getFileType(attrs.ContentType),
				Created:     attrs.Created,
			})
		}

		json.NewEncoder(w).Encode(filesMetaData)
	}
}

func getFileType(contentType string) string {
	parts := strings.Split(contentType, "/")
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}
