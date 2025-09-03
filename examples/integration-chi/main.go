package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/senforsce/tndr"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", tndr.Handler(Home()).ServeHTTP)
	http.ListenAndServe(":3000", r)
}
