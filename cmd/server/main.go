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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"filespace-backend/api/auth"
	user "filespace-backend/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB.")

	var user user.Model
	err = client.Database("test").Collection("users").FindOne(context.Background(), bson.M{"email": "paralejasjustine15@gmail.com"}).Decode(&user)
	if err != nil {
		log.Fatal("Error finding user:", err)
	}

	fmt.Printf("User: %+v\n", user)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Post("/api/v1/auth", auth.Handler(client))
	http.ListenAndServe(os.Getenv("PORT"), router)
}
