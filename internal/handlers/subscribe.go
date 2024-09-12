package handlers

import (
	"encoding/json"
	"gh-checker/internal/config"
	"gh-checker/internal/database"
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

	err := services.UpdateFollowers(req.Followed, config.AppConfig.FollowerUpdateInterval)
	if err != nil {
		respondWithError(w, err)
		return
	}

	isFollowing, err := database.IsFollowing(req.Follower, req.Followed)
	if err != nil {
		respondWithError(w, err)
		return
	}

	json.NewEncoder(w).Encode(models.SubscribeResponse{IsFollowing: isFollowing})
}

func respondWithError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(models.SubscribeResponse{Error: err.Error()})
}
