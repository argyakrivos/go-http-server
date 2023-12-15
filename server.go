package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
)

type Book struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Author string             `json:"author,omitempty" bson:"author,omitempty"`
}

const (
	serverAddress   = "localhost:8080"
	mongoUri        = "mongodb://localhost:27017"
	mongoDb         = "dev"
	mongoCollection = "books"
)

var (
	mongoClient    *mongo.Client
	bookCollection *mongo.Collection
)

func main() {
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	mongoClient = connectToMongoDB(mongoUri)
	bookCollection = mongoClient.Database(mongoDb).Collection(mongoCollection)

	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex).Methods(http.MethodGet)
	r.HandleFunc("/books", handleGetAllBooks).Methods(http.MethodGet)
	r.HandleFunc("/books", handleCreateBook).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", handleGetBook).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", handleUpdateBook).Methods(http.MethodPatch)
	r.HandleFunc("/books/{id}", handleDeleteBook).Methods(http.MethodDelete)
	http.Handle("/", r)

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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	cursor, err := bookCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching books from MongoDB", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var books []Book
	if err := cursor.All(context.Background(), &books); err != nil {
		http.Error(w, "Error decoding books from MongoDB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	result, err := bookCollection.InsertOne(context.Background(), newBook)
	if err != nil {
		http.Error(w, "Error inserting book into MongoDB", http.StatusInternalServerError)
		return
	}
	newBook.Id = result.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func handleGetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	objectId, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	err = bookCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	var updatedBook Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	bookId := vars["id"]

	objectId, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var existingBook Book
	err = bookCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&existingBook)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	setUpdates := bson.M{}
	if updatedBook.Title != "" && updatedBook.Title != existingBook.Title {
		setUpdates["title"] = updatedBook.Title
	}
	if updatedBook.Author != "" && updatedBook.Author != existingBook.Author {
		setUpdates["author"] = updatedBook.Author
	}

	if len(setUpdates) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingBook)
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": setUpdates}
	_, err = bookCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating book in MongoDB", http.StatusInternalServerError)
		return
	}

	bookCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&updatedBook)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	err = bookCollection.FindOneAndDelete(context.Background(), bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
