package services

import (
	"gh-checker/internal/database"
	"gh-checker/internal/lib/logger"
	"time"
)

// UpdateStars проверяет, нужно ли обновить звезды и обновляет их, если необходимо.
// Если обновление не требуется, возвращает кэшированные данные.
func UpdateStars(username, repository string, updateInterval time.Duration) (bool, error) {
	logger.Info("Starting star update process for user " + username + " on repository " + repository)

	// Проверка необходимости обновления звёзд
	shouldUpdate, err := database.ShouldUpdateStars(username, repository, updateInterval)
	if err != nil {
		logger.Error("Error checking if stars need to be updated for user "+username, err)
		return false, err
	}

	if !shouldUpdate {
		logger.Info("No update needed for user " + username + " on repository " + repository)
		return database.IsStarred(username, repository)
	}

	// Обновление звёзд через GitHub API
	hasStar, err := CheckStar(username, repository) // Здесь должен быть вызов GitHub API
	if err != nil {
		logger.Error("Error retrieving stars from GitHub API for user "+username+" on repository "+repository, err)
		return false, err
	}

	// Очистка старых данных о звездах
	err = database.ClearStars(username)
	if err != nil {
		logger.Error("Error clearing stars for user "+username, err)
		return false, err
	}

	// Добавление новых данных о звёздах
	if hasStar {
		err = database.AddStar(username, repository)
		if err != nil {
			logger.Error("Error adding star for user "+username+" on repository "+repository, err)
			return false, err
		}
	}

	// Обновление времени последней проверки звёзд
	err = database.UpdateLastCheckedStars(username, repository) // Используем функцию для звезд
	if err != nil {
		logger.Error("Error updating last checked timestamp for user "+username+" on repository "+repository, err)
		return false, err
	}

	logger.Info("Successfully updated stars for user " + username + " on repository " + repository)
	return hasStar, nil
}
