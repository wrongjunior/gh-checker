package services

import (
	"gh-checker/internal/database"
	"log"
	"time"
)

func UpdateFollowers(username string, updateInterval time.Duration) (bool, error) {
	shouldUpdate, err := database.ShouldUpdateFollowers(username, updateInterval)
	if err != nil {
		return false, err
	}

	if !shouldUpdate {
		return false, nil // Обновление не требуется
	}

	followers, err := GetFollowers(username)
	if err != nil {
		return false, err
	}

	err = database.ClearFollowers(username)
	if err != nil {
		return false, err
	}

	for _, follower := range followers {
		err := database.AddFollower(username, follower)
		if err != nil {
			log.Printf("Error adding follower %s -> %s: %v", follower, username, err)
		}
	}

	err = database.UpdateLastChecked(username)
	if err != nil {
		return false, err
	}

	return true, nil // Обновление выполнено
}
