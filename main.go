package main

import (
	"log"
	"net/http"

	"gh-checker/internal/config"
	"gh-checker/internal/handlers"
	"gh-checker/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig("config.yaml")

	githubAPIKey := config.AppConfig.GitHub.APIKey
	if githubAPIKey == "" {
		log.Fatal("GitHub API key not found in config file")
	}

	services.SetGitHubAPIKey(githubAPIKey)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/subscribe", handlers.SubscribeHandler)

	http.ListenAndServe(":8080", r)
}
