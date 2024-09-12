package services

import (
	"gh-checker/internal/database"
	"log"
	"time"
)

func UpdateFollowers(checkInterval int) {
	for {
		followersToCheck, err := database.GetFollowersToCheck(checkInterval)
		if err != nil {
			log.Printf("Error getting followers to check: %v", err)
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}

		for _, follower := range followersToCheck {
			followers, err := GetFollowers(follower)
			if err != nil {
				log.Printf("Error getting followers for %s: %v", follower, err)
				continue
			}

			for _, followed := range followers {
				err := database.AddFollower(follower, followed)
				if err != nil {
					log.Printf("Error adding follower %s -> %s: %v", follower, followed, err)
				}
			}

			err = database.UpdateLastChecked(follower)
			if err != nil {
				log.Printf("Error updating last checked for %s: %v", follower, err)
			}
		}

		time.Sleep(time.Duration(checkInterval) * time.Second)
	}
}
