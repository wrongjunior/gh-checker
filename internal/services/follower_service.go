package services

import (
	"gh-checker/internal/database"
	"log"
	"time"
)

// UpdateFollowers проверяет, нужно ли обновить подписчиков и обновляет их, если необходимо.
// Если обновление не требуется, возвращает кэшированные данные.
func UpdateFollowers(username string, updateInterval time.Duration) ([]string, bool, error) {
	log.Printf("Checking if we should update followers for %s", username)
	shouldUpdate, err := database.ShouldUpdateFollowers(username, updateInterval)
	if err != nil {
		return nil, false, err
	}

	if !shouldUpdate {
		log.Printf("No need to update followers for %s. Using cached data.", username)
		// Возвращаем кэшированные данные
		followers, err := database.GetFollowers(username)
		if err != nil {
			return nil, false, err
		}
		return followers, false, nil
	}

	log.Printf("Updating followers for %s from GitHub API", username)
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
			log.Printf("Error adding follower %s -> %s: %v", follower, username, err)
		}
	}

	err = database.UpdateLastChecked(username)
	if err != nil {
		return nil, false, err
	}

	log.Printf("Successfully updated followers for %s", username)
	return newFollowers, true, nil
}
