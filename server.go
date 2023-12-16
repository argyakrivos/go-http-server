package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-http-server/db"
	"go-http-server/models"
	"go-http-server/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
)

const (
	serverAddress   = "localhost:8080"
	mongoUri        = "mongodb://localhost:27017"
	mongoDbName     = "dev"
	mongoCollection = "books"
)

type Env struct {
	books *service.DbService
}

func main() {
	mongoClient := connectToMongoDB(mongoUri)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	bookCollection := mongoClient.Database(mongoDbName).Collection(mongoCollection)
	mongoDb := &db.MongoDB{Collection: bookCollection}
	booksService := &service.DbService{DB: mongoDb}
	env := &Env{books: booksService}

	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex).Methods(http.MethodGet)
	r.HandleFunc("/books", env.handleGetAllBooks).Methods(http.MethodGet)
	r.HandleFunc("/books", env.handleCreateBook).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", env.handleGetBook).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", env.handleUpdateBook).Methods(http.MethodPatch)
	r.HandleFunc("/books/{id}", env.handleDeleteBook).Methods(http.MethodDelete)
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

func (e *Env) handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := e.books.GetAllBooks()
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (e *Env) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	resultBook, err := e.books.CreateBook(newBook)
	if err != nil {
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resultBook)
}

func (e *Env) handleGetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	book, err := e.books.GetBook(bookId)
	if err != nil {
		http.Error(w, "Error getting a book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (e *Env) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	var updatedBook models.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	bookId := vars["id"]

	resultBook, err := e.books.UpdateBook(bookId, updatedBook)
	if err != nil {
		http.Error(w, "Error while updating book", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultBook)
}

func (e *Env) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	deletedBook, err := e.books.DeleteBook(bookId)
	if err != nil {
		http.Error(w, "Error while deleting book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedBook)
}
