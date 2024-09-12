package database

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB   *sql.DB
	lock sync.Mutex
)

func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
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
		log.Fatal(err)
	}
}

func AddFollower(username, follower string) error {
	lock.Lock()
	defer lock.Unlock()

	stmt, err := DB.Prepare("INSERT OR IGNORE INTO followers(username, follower, last_updated) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, follower, time.Now())
	return err
}

func IsFollowing(follower, username string) (bool, error) {
	lock.Lock()
	defer lock.Unlock()

	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM followers WHERE username = ? AND follower = ?", username, follower).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateLastChecked(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("INSERT OR REPLACE INTO last_check(username, last_checked) VALUES(?, ?)", username, time.Now())
	return err
}

func ShouldUpdateFollowers(username string, updateInterval time.Duration) (bool, error) {
	lock.Lock()
	defer lock.Unlock()

	var lastChecked time.Time
	err := DB.QueryRow("SELECT last_checked FROM last_check WHERE username = ?", username).Scan(&lastChecked)
	if err == sql.ErrNoRows {
		// Если информации нет в базе, нужно обновить
		return true, nil
	} else if err != nil {
		return false, err
	}
	// Проверяем, истек ли интервал обновления
	return time.Since(lastChecked) > updateInterval, nil
}

func GetFollowers(username string) ([]string, error) {
	lock.Lock()
	defer lock.Unlock()

	rows, err := DB.Query("SELECT follower FROM followers WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []string
	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return followers, nil
}

func ClearFollowers(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("DELETE FROM followers WHERE username = ?", username)
	return err
}
