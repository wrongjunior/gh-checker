package services

import (
	"gh-checker/internal/database"
	"gh-checker/internal/lib/logger"
	"time"
)

// UpdateFollowers проверяет, нужно ли обновить подписчиков и обновляет их, если необходимо.
// Если обновление не требуется, возвращает кэшированные данные.
func UpdateFollowers(username string, updateInterval time.Duration) ([]string, bool, error) {
	logger.Info("Checking if we should update followers for " + username)
	shouldUpdate, err := database.ShouldUpdateFollowers(username, updateInterval)
	if err != nil {
		return nil, false, err
	}

	if !shouldUpdate {
		logger.Info("No need to update followers for " + username + ". Using cached data.")
		// Возвращаем кэшированные данные
		followers, err := database.GetFollowers(username)
		if err != nil {
			return nil, false, err
		}
		return followers, false, nil
	}

	logger.Info("Updating followers for " + username + " from GitHub API")
	newFollowers, err := GetFollowers(username)
	if err != nil {
		return nil, false, err
	}

	err = database.ClearFollowers(username)
	if err != nil {
		return nil, false, err
	}

	for _, follower := range newFollowers {
		err := database.AddFollower(username, follower)
		if err != nil {
			logger.Error("Error adding follower "+follower+" -> "+username, err)
		}
	}

	err = database.UpdateLastChecked(username)
	if err != nil {
		return nil, false, err
	}

	logger.Info("Successfully updated followers for " + username)
	return newFollowers, true, nil
}
