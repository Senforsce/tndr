package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/senforsce/tndr"
)

func main() {
	mux := http.NewServeMux()

	// Serve the tndr page.
	mux.Handle("/", tndr.Handler(page()))

	// Serve static content.
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server.
	fmt.Println("listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}
