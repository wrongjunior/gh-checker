package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

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
		follower TEXT NOT NULL,
		followed TEXT NOT NULL,
		last_checked TIMESTAMP,
		UNIQUE(follower, followed)
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func AddFollower(follower, followed string) error {
	stmt, err := DB.Prepare("INSERT OR REPLACE INTO followers(follower, followed, last_checked) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(follower, followed, time.Now())
	return err
}

func GetFollowersToCheck(checkInterval int) ([]string, error) {
	query := `
	SELECT DISTINCT follower 
	FROM followers 
	WHERE last_checked IS NULL OR last_checked < datetime('now', '-' || ? || ' seconds')`

	rows, err := DB.Query(query, checkInterval)
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

func UpdateLastChecked(follower string) error {
	_, err := DB.Exec("UPDATE followers SET last_checked = ? WHERE follower = ?", time.Now(), follower)
	return err
}

func IsFollowing(follower, followed string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? AND followed = ?", follower, followed).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
