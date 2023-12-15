package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	books []Book
)

func main() {
	books = []Book{
		{Id: 1, Title: "Go Programming", Author: "John Doe"},
		{Id: 2, Title: "Web Development with Go", Author: "Jane Smith"},
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex).Methods(http.MethodGet)
	r.HandleFunc("/books", handleGetAllBooks).Methods(http.MethodGet)
	r.HandleFunc("/books", handleCreateBook).Methods(http.MethodPost)
	r.HandleFunc("/books/{id:[0-9]+}", handleGetBook).Methods(http.MethodGet)
	r.HandleFunc("/books/{id:[0-9]+}", handleUpdateBook).Methods(http.MethodPatch)
	r.HandleFunc("/books/{id:[0-9]+}", handleDeleteBook).Methods(http.MethodDelete)
	http.Handle("/", r)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
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

	newBook.Id = len(books) + 1
	books = append(books, newBook)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func handleGetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, _ := strconv.Atoi(vars["id"])
	for _, book := range books {
		if book.Id == bookId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.Error(w, "book not found", http.StatusNotFound)
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	var updatedBook Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	bookId, _ := strconv.Atoi(vars["id"])
	for i, book := range books {
		if book.Id == bookId {
			books[i].Title = updatedBook.Title
			books[i].Author = updatedBook.Author

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(books[i])
			return
		}
	}

	http.Error(w, "book not found", http.StatusNotFound)
}

func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, _ := strconv.Atoi(vars["id"])
	for i, book := range books {
		if book.Id == bookId {
			books = append(books[:i], books[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.Error(w, "book not found", http.StatusNotFound)
}
