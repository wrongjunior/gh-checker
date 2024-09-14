package handlers

import (
	"encoding/json"
	"gh-checker/internal/config"
	"gh-checker/internal/lib/logger" // Импортируем твой логгер
	"gh-checker/internal/models"
	"gh-checker/internal/services"
	"net/http"
)

// respondWithError отвечает с ошибкой и логирует её
func respondWithError(w http.ResponseWriter, err error) {
	logger.Error("Responding with error", err)
	w.WriteHeader(http.StatusInternalServerError)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encodeErr != nil {
		logger.Error("Error while encoding error response", encodeErr) // Логируем ошибку кодирования
	}
}

// SubscribeHandler обрабатывает запрос на проверку подписчиков
func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Processing SubscribeHandler request")
	var req models.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Invalid request body", err)
		return
	}
	logger.Info("Request body successfully decoded")

	logger.Info("Received request to check if " + req.Follower + " is following " + req.Followed)

	logger.Info("Calling UpdateFollowers service for user " + req.Followed)
	followers, updated, err := services.UpdateFollowers(req.Followed, config.AppConfig.FollowerUpdateInterval)
	if err != nil {
		logger.Error("Error while updating followers", err)
		respondWithError(w, err)
		return
	}
	logger.Info("UpdateFollowers service call succeeded")

	isFollowing := false
	logger.Info("Checking if " + req.Follower + " is in the list of followers")
	for _, follower := range followers {
		if follower == req.Follower {
			isFollowing = true
			logger.Info(req.Follower + " is following " + req.Followed)
			break
		}
	}
	if !isFollowing {
		logger.Info(req.Follower + " is not following " + req.Followed)
	}

	if updated {
		logger.Info("Followers list for " + req.Followed + " was updated from GitHub API")
	} else {
		logger.Info("Using cached followers data for " + req.Followed)
	}

	response := models.SubscribeResponse{IsFollowing: isFollowing}
	logger.Info("Encoding response to JSON")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Error while encoding response to JSON", err)
		respondWithError(w, err)
		return
	}
	logger.Info("Response successfully sent to the client")
}
