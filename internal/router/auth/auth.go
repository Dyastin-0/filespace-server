package router

import (
	"github.com/go-chi/chi/v5"

	"go.mongodb.org/mongo-driver/mongo"

	auth "filespace/internal/auth"
)

func Auth(client *mongo.Client) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", auth.Handler(client))
	router.Post("/sign-up", auth.Signup(client))
	router.Post("/refresh", auth.Refresh(client))
	router.Post("/verify", auth.Verify(client))

	return router
}
