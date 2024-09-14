package services

import (
	"gh-checker/internal/database"
	"gh-checker/internal/lib/logger"
	"time"
)

// UpdateFollowers проверяет, нужно ли обновить подписчиков и обновляет их, если необходимо.
// Если обновление не требуется, возвращает кэшированные данные.
func UpdateFollowers(username string, updateInterval time.Duration) ([]string, bool, error) {
	logger.Info("Starting follower update process for user " + username)

	// Проверка необходимости обновления подписчиков
	logger.Info("Checking if followers need to be updated for user " + username)
	shouldUpdate, err := database.ShouldUpdateFollowers(username, updateInterval)
	if err != nil {
		logger.Error("Error checking if followers need to be updated for user "+username, err)
		return nil, false, err
	}

	if !shouldUpdate {
		logger.Info("No update needed for user " + username + ". Retrieving cached followers.")
		// Возвращаем кэшированные данные
		followers, err := database.GetFollowers(username)
		if err != nil {
			logger.Error("Error retrieving cached followers for user "+username, err)
			return nil, false, err
		}
		logger.Info("Successfully retrieved cached followers for user " + username)
		return followers, false, nil
	}

	// Обновление подписчиков через GitHub API
	logger.Info("Updating followers for user " + username + " via GitHub API")
	newFollowers, err := GetFollowers(username) // Здесь должен быть вызов GitHub API
	if err != nil {
		logger.Error("Error retrieving followers from GitHub API for user "+username, err)
		return nil, false, err
	}

	// Очистка старых подписчиков
	logger.Info("Clearing old followers for user " + username)
	err = database.ClearFollowers(username)
	if err != nil {
		logger.Error("Error clearing followers for user "+username, err)
		return nil, false, err
	}

	// Добавление новых подписчиков в базу данных
	logger.Info("Adding new followers for user " + username)
	for _, follower := range newFollowers {
		logger.Info("Adding follower " + follower + " for user " + username)
		err := database.AddFollower(username, follower)
		if err != nil {
			logger.Error("Error adding follower "+follower+" -> "+username, err)
		} else {
			logger.Info("Successfully added follower " + follower + " for user " + username)
		}
	}

	// Обновление времени последней проверки подписчиков
	logger.Info("Updating last checked timestamp for user " + username)
	err = database.UpdateLastCheckedFollowers(username) // Используем функцию для подписчиков
	if err != nil {
		logger.Error("Error updating last checked timestamp for user "+username, err)
		return nil, false, err
	}

	logger.Info("Successfully updated followers for user " + username)
	return newFollowers, true, nil
}
