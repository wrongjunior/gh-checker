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
	log.Printf("Responding with error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Invalid request body")
		return
	}

	log.Printf("Received request to check if %s is following %s", req.Follower, req.Followed)

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
		log.Printf("Updated followers for %s from GitHub API", req.Followed)
	} else {
		log.Printf("Used cached followers data for %s", req.Followed)
	}

	json.NewEncoder(w).Encode(models.SubscribeResponse{IsFollowing: isFollowing})
}
