package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!")
	})
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
