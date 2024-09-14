package handlers

import (
	"encoding/json"
	"fmt"
	"gh-checker/internal/lib/logger"
	"gh-checker/internal/models"
	"gh-checker/internal/services"
	"net/http"
)

// StarCheckHandler обрабатывает запрос на проверку, поставил ли пользователь звезду на репозиторий
func StarCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Processing StarCheckHandler request")

	var req models.StarCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Invalid request body", err)
		return
	}

	logger.Info(fmt.Sprintf("Received request to check if %s starred repository %s", req.Username, req.Repository))

	hasStar, err := services.CheckStar(req.Username, req.Repository)
	if err != nil {
		logger.Error("Error while checking star", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.StarCheckResponse{HasStar: hasStar}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Error encoding response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info(fmt.Sprintf("Successfully responded to star check for user %s on repository %s", req.Username, req.Repository))
}
