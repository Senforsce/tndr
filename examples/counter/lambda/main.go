package main

import (
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/senforsce/tndr/examples/counter/db"
	"github.com/senforsce/tndr/examples/counter/handlers"
	"github.com/senforsce/tndr/examples/counter/services"
	"github.com/senforsce/tndr/examples/counter/session"
	"golang.org/x/exp/slog"
)

func main() {
	// Create handlers.
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	s, err := db.NewCountStore(os.Getenv("TABLE_NAME"), os.Getenv("AWS_REGION"))
	if err != nil {
		log.Error("failed to create store", slog.Any("error", err))
		os.Exit(1)
	}
	cs := services.NewCount(log, s)
	h := handlers.New(log, cs)

	// Add session middleware.
	sh := session.NewMiddleware(h)

	// Start Lambda.
	algnhsa.ListenAndServe(sh, nil)
}
