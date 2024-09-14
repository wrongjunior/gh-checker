package database

import (
	"database/sql"
	"gh-checker/internal/lib/logger"
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
    username TEXT NOT NULL,
    repository TEXT NOT NULL,
    last_checked TIMESTAMP,
    UNIQUE(username, repository)
	);
	CREATE TABLE IF NOT EXISTS stars (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		repository TEXT NOT NULL,
		last_updated TIMESTAMP,
		UNIQUE(username, repository)
	);
	`

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

// UpdateLastChecked обновляет время последней проверки подписчиков для пользователя
func UpdateLastChecked(username, recordType string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("INSERT OR REPLACE INTO last_check(username, repository, last_checked) VALUES(?, ?, ?)", username, recordType, time.Now())
	if err != nil {
		logger.Error("Error updating last checked time for user and record type", err)
		return err
	}

	logger.Info("Updated last checked time for user " + username + " and record type " + recordType)
	return nil
}

// UpdateLastCheckedFollowers обновляет время последней проверки подписчиков для пользователя
func UpdateLastCheckedFollowers(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("INSERT OR REPLACE INTO last_check(username, repository, last_checked) VALUES(?, ?, ?)", username, "followers", time.Now())
	if err != nil {
		logger.Error("Error updating last checked time for user and followers", err)
		return err
	}

	logger.Info("Updated last checked time for user " + username + " for followers")
	return nil
}

// UpdateLastCheckedStars обновляет время последней проверки звезд для пользователя и репозитория
func UpdateLastCheckedStars(username, repository string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("INSERT OR REPLACE INTO last_check(username, repository, last_checked) VALUES(?, ?, ?)", username, repository, time.Now())
	if err != nil {
		logger.Error("Error updating last checked time for user and repository", err)
		return err
	}

	logger.Info("Updated last checked time for user " + username + " and repository " + repository)
	return nil
}

func ShouldUpdateFollowers(username string, updateInterval time.Duration) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	var lastChecked time.Time
	err := DB.QueryRow("SELECT last_checked FROM last_check WHERE username = ? AND repository = 'followers'", username).Scan(&lastChecked)
	if err == sql.ErrNoRows {
		logger.Info("No last checked time found for user " + username + ". Update required.")
		return true, nil
	} else if err != nil {
		logger.Error("Error checking last checked time for user", err)
		return false, err
	}

	timeSinceLastCheck := time.Since(lastChecked)
	return timeSinceLastCheck > updateInterval, nil
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

// AddStar добавляет информацию о звезде пользователя на репозитории
func AddStar(username, repository string) error {
	lock.Lock()
	defer lock.Unlock()

	stmt, err := DB.Prepare("INSERT OR IGNORE INTO stars(username, repository, last_updated) VALUES(?, ?, ?)")
	if err != nil {
		logger.Error("Error preparing statement for adding star", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, repository, time.Now())
	if err != nil {
		logger.Error("Error executing statement for adding star", err)
		return err
	}

	logger.Info("Added/updated star for user " + username + " on repository " + repository)
	return nil
}

// IsStarred проверяет, поставил ли пользователь звезду на репозиторий
func IsStarred(username, repository string) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM stars WHERE username = ? AND repository = ?", username, repository).Scan(&count)
	if err != nil {
		logger.Error("Error checking if user starred repository", err)
		return false, err
	}

	logger.Info("Checked if user " + username + " starred repository " + repository)
	return count > 0, nil
}

// ClearStars удаляет все звезды пользователя на репозитории
func ClearStars(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("DELETE FROM stars WHERE username = ?", username)
	if err != nil {
		logger.Error("Error clearing stars for user", err)
		return err
	}

	logger.Info("Cleared stars for user " + username)
	return nil
}

// GetLastChecked возвращает время последней проверки для пользователя и репозитория
func GetLastChecked(username, repository string) (time.Time, error) {
	lock.RLock()
	defer lock.RUnlock()

	var lastChecked time.Time
	err := DB.QueryRow("SELECT last_checked FROM last_check WHERE username = ? AND repository = ?", username, repository).Scan(&lastChecked)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("No last checked time found for user " + username + " and repository " + repository)
			return time.Time{}, sql.ErrNoRows
		}
		logger.Error("Error retrieving last checked time for user and repository", err)
		return time.Time{}, err
	}

	logger.Info("Retrieved last checked time for user " + username + " and repository " + repository)
	return lastChecked, nil
}

func ShouldUpdateStars(username, repository string, updateInterval time.Duration) (bool, error) {
	lock.RLock()
	defer lock.RUnlock()

	var lastChecked time.Time
	err := DB.QueryRow("SELECT last_checked FROM last_check WHERE username = ? AND repository = ?", username, repository).Scan(&lastChecked)
	if err == sql.ErrNoRows {
		logger.Info("No last checked time found for user " + username + " and repository " + repository + ". Update required.")
		return true, nil
	} else if err != nil {
		logger.Error("Error checking last checked time for user and repository", err)
		return false, err
	}

	timeSinceLastCheck := time.Since(lastChecked)
	return timeSinceLastCheck > updateInterval, nil
}
