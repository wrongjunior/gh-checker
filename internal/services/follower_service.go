package services

import (
	"gh-checker/internal/database"
	"log"
	"time"
)

func UpdateFollowers(username string, updateInterval time.Duration) error {
	shouldUpdate, err := database.ShouldUpdateFollowers(username, updateInterval)
	if err != nil {
		return err
	}

	if !shouldUpdate {
		return nil
	}

	followers, err := GetFollowers(username)
	if err != nil {
		return err
	}

	err = database.ClearFollowers(username)
	if err != nil {
		return err
	}

	for _, follower := range followers {
		err := database.AddFollower(username, follower)
		if err != nil {
			log.Printf("Error adding follower %s -> %s: %v", follower, username, err)
		}
	}

	return database.UpdateLastChecked(username)
}
