package services

import (
	"gh-checker/internal/database"
	"log"
	"time"
)

// UpdateFollowers проверяет, нужно ли обновить подписчиков и обновляет их, если необходимо.
// Если обновление не требуется, возвращает кэшированные данные.
func UpdateFollowers(username string, updateInterval time.Duration) ([]string, bool, error) {
	// Проверяем, нужно ли обновить подписчиков
	shouldUpdate, err := database.ShouldUpdateFollowers(username, updateInterval)
	if err != nil {
		return nil, false, err
	}

	if !shouldUpdate {
		// Возвращаем кэшированные данные
		followers, err := database.GetFollowers(username)
		if err != nil {
			return nil, false, err
		}
		return followers, false, nil
	}

	// Если нужно обновить подписчиков с GitHub
	newFollowers, err := GetFollowers(username)
	if err != nil {
		return nil, false, err
	}

	// Очищаем текущих подписчиков и добавляем новых
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

	// Обновляем время последней проверки
	err = database.UpdateLastChecked(username)
	if err != nil {
		return nil, false, err
	}

	return newFollowers, true, nil
}
