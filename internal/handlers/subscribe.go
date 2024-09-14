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
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// SubscribeHandler обрабатывает запрос на проверку подписчиков
func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Invalid request body", err)
		return
	}

	logger.Info("Received request to check if " + req.Follower + " is following " + req.Followed)

	followers, updated, err := services.UpdateFollowers(req.Followed, config.AppConfig.FollowerUpdateInterval)
	if err != nil {
		respondWithError(w, err)
		return
	}

	isFollowing := false
	for _, follower := range followers {
		if follower == req.Follower {
			isFollowing = true
			break
		}
	}

	if updated {
		logger.Info("Updated followers for " + req.Followed + " from GitHub API")
	} else {
		logger.Info("Used cached followers data for " + req.Followed)
	}

	if err := json.NewEncoder(w).Encode(models.SubscribeResponse{IsFollowing: isFollowing}); err != nil {
		respondWithError(w, err)
		return
	}
}
