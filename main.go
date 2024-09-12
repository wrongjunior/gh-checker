package main

import (
	"gh-checker/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Ручка для подписки
	r.Post("/api/subscribe", handlers.SubscribeHandler)

	http.ListenAndServe(":8080", r)
}
