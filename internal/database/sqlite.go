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
	log.Printf("Initialized database connection to %s", dbPath)

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
	log.Println("Created necessary tables if they did not exist")
}

func AddFollower(username, follower string) error {
	lock.Lock()
	defer lock.Unlock()

	stmt, err := DB.Prepare("INSERT OR IGNORE INTO followers(username, follower, last_updated) VALUES(?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement for adding follower %s -> %s: %v", follower, username, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, follower, time.Now())
	if err != nil {
		log.Printf("Error executing statement for adding follower %s -> %s: %v", follower, username, err)
	}
	log.Printf("Added/updated follower %s for user %s", follower, username)
	return err
}

func IsFollowing(follower, username string) (bool, error) {
	lock.Lock()
	defer lock.Unlock()

	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM followers WHERE username = ? AND follower = ?", username, follower).Scan(&count)
	if err != nil {
		log.Printf("Error checking if follower %s follows %s: %v", follower, username, err)
		return false, err
	}
	log.Printf("Checked if follower %s follows %s: %v", follower, username, count > 0)
	return count > 0, nil
}

func UpdateLastChecked(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("INSERT OR REPLACE INTO last_check(username, last_checked) VALUES(?, ?)", username, time.Now())
	if err != nil {
		log.Printf("Error updating last checked time for %s: %v", username, err)
	} else {
		log.Printf("Updated last checked time for %s", username)
	}
	return err
}

func ShouldUpdateFollowers(username string, updateInterval time.Duration) (bool, error) {
	lock.Lock()
	defer lock.Unlock()

	var lastChecked time.Time
	err := DB.QueryRow("SELECT last_checked FROM last_check WHERE username = ?", username).Scan(&lastChecked)
	if err == sql.ErrNoRows {
		log.Printf("No last checked time found for %s. Update required.", username)
		return true, nil
	} else if err != nil {
		log.Printf("Error checking last checked time for %s: %v", username, err)
		return false, err
	}

	timeSinceLastCheck := time.Since(lastChecked)
	log.Printf("Time since last checked for %s: %v", username, timeSinceLastCheck)

	shouldUpdate := timeSinceLastCheck > updateInterval
	log.Printf("Should update followers for %s: %v", username, shouldUpdate)
	return shouldUpdate, nil
}

func GetFollowers(username string) ([]string, error) {
	lock.Lock()
	defer lock.Unlock()

	rows, err := DB.Query("SELECT follower FROM followers WHERE username = ?", username)
	if err != nil {
		log.Printf("Error retrieving followers for %s: %v", username, err)
		return nil, err
	}
	defer rows.Close()

	var followers []string
	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			log.Printf("Error scanning follower for %s: %v", username, err)
			return nil, err
		}
		followers = append(followers, follower)
	}

	log.Printf("Retrieved %d followers for %s", len(followers), username)
	return followers, nil
}

func ClearFollowers(username string) error {
	lock.Lock()
	defer lock.Unlock()

	_, err := DB.Exec("DELETE FROM followers WHERE username = ?", username)
	if err != nil {
		log.Printf("Error clearing followers for %s: %v", username, err)
		return err
	}
	log.Printf("Cleared followers for %s", username)
	return nil
}
