package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	books          []Book
	IndexExp       = regexp.MustCompile(`^/?$`)
	GetAllBooksExp = regexp.MustCompile(`^/books/?$`)
	GetBookExp     = regexp.MustCompile(`^/books/(\d+)?$`)
)

func main() {
	books = []Book{
		{Id: 1, Title: "Go Programming", Author: "John Doe"},
		{Id: 2, Title: "Web Development with Go", Author: "Jane Smith"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && IndexExp.MatchString(r.URL.Path):
			handleIndex(w, r)
			return
		default:
			http.NotFound(w, r)
			return
		}
	})

	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		getBookMatch := GetBookExp.FindStringSubmatch(r.URL.Path)
		switch {
		case r.Method == http.MethodGet && GetAllBooksExp.MatchString(r.URL.Path):
			handleGetAllBooks(w, r)
			return
		case r.Method == http.MethodGet && getBookMatch != nil:
			bookId, _ := strconv.Atoi(getBookMatch[1])
			handleGetBook(w, r, bookId)
			return
		default:
			http.NotFound(w, r)
			return
		}
	})

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

func handleGetBook(w http.ResponseWriter, r *http.Request, bookId int) {
	for _, book := range books {
		if book.Id == bookId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, r)
}
