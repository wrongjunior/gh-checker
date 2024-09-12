package handlers

import (
	"encoding/json"
	"gh-checker/internal/models"
	"gh-checker/internal/services"
	"net/http"
)

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для проверки подписки
	isFollowing, err := services.CheckIfFollowing(req.Follower, req.Followed)
	if err != nil {
		json.NewEncoder(w).Encode(models.SubscribeResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(models.SubscribeResponse{IsFollowing: isFollowing})
}
