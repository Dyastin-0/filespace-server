package router

import (
	"github.com/go-chi/chi/v5"

	"go.mongodb.org/mongo-driver/mongo"

	auth "filespace/internal/auth"
)

func Auth(client *mongo.Client) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", auth.Handler(client))
	router.Post("/signup", auth.Signup(client))
	router.Get("/refresh", auth.Refresh(client))

	return router
}
