package main

import (
	"log"
	"net/http"

	"casino_backend/internal/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Регистрируем API handlers
	api.NewAuthHandler(r, api.AuthHandlerDeps{})
	api.NewLineHandler(r, api.LineHandlerDeps{})
	api.NewUserHandler(r, api.UserHandlerDeps{})

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
