package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cloud.google.com/go/storage"

	middleware "filespace/internal/middleware"
	router "filespace/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	version := os.Getenv("VERSION")

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}

	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal("Failed to create storage client: ", err)
	}

	fmt.Println("Connected to MongoDB.")

	MainRouter := chi.NewRouter()

	MainRouter.Use(middleware.Logger)
	MainRouter.Use(middleware.Credential)
	MainRouter.Use(render.SetContentType(render.ContentTypeJSON))

	MainRouter.Mount("/api/"+version+"/auth", router.Auth(mongoClient))
	MainRouter.Mount("/api/"+version+"/files", router.File(storageClient, mongoClient))

	port := os.Getenv("PORT")
	if err := http.ListenAndServe(":"+port, MainRouter); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
