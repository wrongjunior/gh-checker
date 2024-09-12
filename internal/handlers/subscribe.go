package handlers

import (
	"encoding/json"
	"gh-checker/internal/config"
	"gh-checker/internal/models"
	"gh-checker/internal/services"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем подписчиков и обновляем их при необходимости
	followers, updated, err := services.UpdateFollowers(req.Followed, config.AppConfig.FollowerUpdateInterval)
	if err != nil {
		respondWithError(w, err)
		return
	}

	// Проверяем, подписан ли `Follower` на `Followed`
	isFollowing := false
	for _, follower := range followers {
		if follower == req.Follower {
			isFollowing = true
			break
		}
	}

	if updated {
		log.Printf("Updated followers for %s from GitHub API", req.Followed)
	} else {
		log.Printf("Used cached followers data for %s", req.Followed)
	}

	// Возвращаем результат
	json.NewEncoder(w).Encode(models.SubscribeResponse{IsFollowing: isFollowing})
}
