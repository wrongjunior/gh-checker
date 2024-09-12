package main

import (
	"log"
	"net/http"

	"gh-checker/internal/config"
	"gh-checker/internal/database"
	"gh-checker/internal/handlers"
	"gh-checker/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.LoadConfig("config.yaml")

	githubAPIKey := config.AppConfig.GitHub.APIKey
	if githubAPIKey == "" {
		log.Fatal("GitHub API key not found in config file")
	}

	services.SetGitHubAPIKey(githubAPIKey)

	database.InitDB(config.AppConfig.Database.Path)

	go services.UpdateFollowers(config.AppConfig.FollowerCheckInterval)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/subscribe", handlers.SubscribeHandler)

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
