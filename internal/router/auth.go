package router

import (
	"github.com/go-chi/chi/v5"

	"go.mongodb.org/mongo-driver/mongo"

	auth "filespace/internal/auth"
	refresh "filespace/internal/refresh"
)

func Auth(client *mongo.Client) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", auth.Handler(client))
	router.Get("/refresh", refresh.Handler(client))

	return router
}
