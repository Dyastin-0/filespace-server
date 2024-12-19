package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	credential "filespace/internal/credential"
	router "filespace/internal/router/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	version := os.Getenv("VERSION")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB.")

	MainRouter := chi.NewRouter()

	MainRouter.Use(middleware.Logger)
	MainRouter.Use(credential.Handler)
	MainRouter.Use(render.SetContentType(render.ContentTypeJSON))

	MainRouter.Mount("/api/"+version+"/auth", router.Auth(client))

	port := os.Getenv("PORT")
	if err := http.ListenAndServe(":"+port, MainRouter); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
