package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/senforsce/tndr"
)

func main() {
	component := hello("John")

	http.Handle("/", tndr.Handler(component))

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
