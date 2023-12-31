package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go-microservice/db"
	"go-microservice/models"
	"go-microservice/service"
	"io"
	"log/slog"
	"net/http"
)

type BookAPI struct {
	bookService service.BookService
}

func (api *BookAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", api.handleIndex).Methods(http.MethodGet)
	router.HandleFunc("/books", api.handleGetAllBooks).Methods(http.MethodGet)
	router.HandleFunc("/books", api.handleCreateBook).Methods(http.MethodPost)
	router.HandleFunc("/books/{id}", api.handleGetBook).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", api.handleUpdateBook).Methods(http.MethodPatch)
	router.HandleFunc("/books/{id}", api.handleDeleteBook).Methods(http.MethodDelete)
}

func NewBookAPI(bookService service.BookService) *BookAPI {
	return &BookAPI{bookService: bookService}
}

func (api *BookAPI) handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func (api *BookAPI) handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := api.bookService.GetAllBooks()
	if err != nil {
		slog.Error("Error fetching books: ", err)
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (api *BookAPI) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	resultBook, err := api.bookService.CreateBook(newBook)
	if err != nil {
		slog.Error("Error creating book: ", err)
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resultBook)
}

func (api *BookAPI) handleGetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	book, err := api.bookService.GetBook(bookId)
	if err != nil {
		if errors.Is(err, db.ErrInvalidID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, db.ErrBookNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			slog.Error("Error getting a book: ", err)
			http.Error(w, "Error getting a book", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (api *BookAPI) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	var updatedBook models.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	bookId := vars["id"]

	resultBook, err := api.bookService.UpdateBook(bookId, updatedBook)
	if err != nil {
		if errors.Is(err, db.ErrInvalidID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, db.ErrBookNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			slog.Error("Error updating a book: ", err)
			http.Error(w, "Error updating book", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultBook)
}

func (api *BookAPI) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	deletedBook, err := api.bookService.DeleteBook(bookId)
	if err != nil {
		if errors.Is(err, db.ErrInvalidID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, db.ErrBookNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			slog.Error("Error deleting a book: ", err)
			http.Error(w, "Error deleting book", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedBook)
}
