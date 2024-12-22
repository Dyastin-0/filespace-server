package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"cloud.google.com/go/storage"

	file "filespace/internal/file"
	"filespace/internal/middleware"
)

func File(storageClient *storage.Client, mongoClient *mongo.Client) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.JWT)
	router.Get("/", file.Get(storageClient))
	router.Post("/", middleware.CheckStorageCapacity(mongoClient)(file.Post(storageClient, mongoClient)))
	router.Delete("/", file.Delete(storageClient, mongoClient))

	return router
}
