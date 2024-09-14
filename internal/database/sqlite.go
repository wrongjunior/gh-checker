package database

import (
	"database/sql"
	"fmt"
	"gh-checker/internal/lib/logger"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB   *sql.DB
	lock sync.RWMutex // Используем RWMutex для разделения блокировки
)

// InitDB инициализирует базу данных
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error("Failed to open database", err)
		return err
	}

	logger.Info("Initialized database connection to " + dbPath)

	if err = createTables(); err != nil {
		return err
	}

	return nil
}

// createTables создаёт необходимые таблицы, если они не существуют
func createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS followers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		follower TEXT NOT NULL,
		last_updated TIMESTAMP,
		UNIQUE(username, follower)
	);
	CREATE TABLE IF NOT EXISTS last_check (
		username TEXT PRIMARY KEY,
		last_checked TIMESTAMP
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		logger.Error("Failed to create tables", err)
		return err
	}

	logger.Info("Created necessary tables")
	return nil
}

// AddFollower добавляет нового подписчика
func AddFollower(username, follower string) error {
	lock.Lock()
	defer lock.Unlock()

	stmt, err := DB.Prepare("INSERT OR IGNORE INTO followers(username, follower, last_updated) VALUES(?, ?, ?)")
	if err != nil {
		logger.Error("Error preparing statement for adding follower", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, follower, time.Now())
	if err != nil {
		logger.Error("Error executing statement for adding follower", err)
		return err
	}

	logger.Info("Added/updated follower " + follower + " for user " + username)
	return nil
}

// IsFollowing проверяет, является ли follower подписчиком username
func IsFollowing(follower, username string) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM followers WHERE username = ? AND follower = ?", username, follower).Scan(&count)
	if err != nil {
		logger.Error("Error checking if follower follows user", err)
		return false, err
	}

	logger.Info("Checked if follower " + follower + " follows user " + username)
	return count > 0, nil
}

// UpdateLastChecked обновляет время последней проверки подписчиков для username
func UpdateLastChecked(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("INSERT OR REPLACE INTO last_check(username, last_checked) VALUES(?, ?)", username, time.Now())
	if err != nil {
		logger.Error("Error updating last checked time for user", err)
		return err
	}

	logger.Info("Updated last checked time for user " + username)
	return nil
}

func ShouldUpdateFollowers(username string, updateInterval time.Duration) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	var lastChecked time.Time
	err := DB.QueryRow("SELECT last_checked FROM last_check WHERE username = ?", username).Scan(&lastChecked)
	if err == sql.ErrNoRows {
		logger.Info("No last checked time found for user " + username + ". Update required.")
		return true, nil
	} else if err != nil {
		logger.Error("Error checking last checked time for user", err)
		return false, err
	}

	timeSinceLastCheck := time.Since(lastChecked)
	logger.Info("Time since last checked for user " + username + ": " + timeSinceLastCheck.String())

	shouldUpdate := timeSinceLastCheck > updateInterval
	logger.Info(fmt.Sprintf("Should update followers for user %s: %s", username, strconv.FormatBool(shouldUpdate)))
	return shouldUpdate, nil
}

// GetFollowers возвращает список подписчиков пользователя
func GetFollowers(username string) ([]string, error) {
	lock.RLock()
	defer lock.RUnlock()

	rows, err := DB.Query("SELECT follower FROM followers WHERE username = ?", username)
	if err != nil {
		logger.Error("Error retrieving followers for user", err)
		return nil, err
	}
	defer rows.Close()

	var followers []string
	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			logger.Error("Error scanning follower for user", err)
			return nil, err
		}
		followers = append(followers, follower)
	}

	logger.Info("Retrieved " + string(len(followers)) + " followers for user " + username)
	return followers, nil
}

// ClearFollowers удаляет всех подписчиков пользователя
func ClearFollowers(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("DELETE FROM followers WHERE username = ?", username)
	if err != nil {
		logger.Error("Error clearing followers for user", err)
		return err
	}

	logger.Info("Cleared followers for user " + username)
	return nil
}
