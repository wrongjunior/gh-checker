package handlers

import (
	"encoding/json"
	"gh-checker/internal/config"
	"gh-checker/internal/lib/logger"
	"gh-checker/internal/models"
	"gh-checker/internal/services"
	"net/http"
)

// respondWithJSON отвечает клиенту с JSON-ответом и заголовком Content-Type
func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error encoding JSON response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// StarCheckHandler обрабатывает запрос на проверку, поставил ли пользователь звезду на репозиторий
func StarCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Processing StarCheckHandler request")

	var req models.StarCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Invalid request body", err)
		return
	}

	logger.Info("Received request to check if " + req.Username + " starred repository " + req.Repository)

	hasStar, err := services.UpdateStars(req.Username, req.Repository, config.AppConfig.FollowerUpdateInterval)
	if err != nil {
		logger.Error("Error while updating stars", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.StarCheckResponse{HasStar: hasStar}

	// Устанавливаем заголовок Content-Type и отвечаем клиенту
	respondWithJSON(w, response)

	logger.Info("Successfully responded to star check for user " + req.Username + " on repository " + req.Repository)
}
