package services

import (
	"gh-checker/internal/database"
	"log"
	"time"
)

func UpdateFollowers(username string, updateInterval time.Duration) ([]string, bool, error) {
	shouldUpdate, err := database.ShouldUpdateFollowers(username, updateInterval)
	if err != nil {
		return nil, false, err
	}

	if !shouldUpdate {
		// Получаем подписчиков из локальной базы данных
		followers, err := database.GetFollowers(username)
		return followers, false, err
	}

	// Если нужно обновить подписчиков с GitHub
	followers, err := GetFollowers(username)
	if err != nil {
		return nil, false, err
	}

	err = database.ClearFollowers(username)
	if err != nil {
		return nil, false, err
	}

	for _, follower := range followers {
		err := database.AddFollower(username, follower)
		if err != nil {
			log.Printf("Error adding follower %s -> %s: %v", follower, username, err)
		}
	}

	err = database.UpdateLastChecked(username)
	if err != nil {
		return nil, false, err
	}

	return followers, true, nil
}
