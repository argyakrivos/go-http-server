package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go-http-server/api"
	"go-http-server/db"
	"go-http-server/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

const (
	serverAddress   = "localhost:8080"
	mongoUri        = "mongodb://localhost:27017"
	mongoDbName     = "dev"
	mongoCollection = "books"
)

func main() {
	// Setup MongoDB connection and get a collection
	mongoClient := connectToMongoDB(mongoUri)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	bookCollection := mongoClient.Database(mongoDbName).Collection(mongoCollection)

	// Initialise the book repository and service
	bookRepository := db.NewMongoDBBookRepository(bookCollection)
	bookService := service.NewBookServiceImpl(bookRepository)

	// Initialise the router
	router := mux.NewRouter()
	http.Handle("/", router)

	// Initialise the API and register routes
	bookAPI := api.NewBookAPI(bookService)
	bookAPI.RegisterRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Server listening on %v", serverAddress))
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func connectToMongoDB(uri string) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	return client
}
