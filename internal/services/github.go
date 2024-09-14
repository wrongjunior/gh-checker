package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const githubAPI = "https://api.github.com"

// Ограничение на количество подписчиков, загружаемых за один запрос
const maxFollowersPerPage = 100

var githubAPIKey string

// SetGitHubAPIKey устанавливает API ключ GitHub
func SetGitHubAPIKey(apiKey string) {
	githubAPIKey = apiKey
	log.Println("GitHub API key set")
}

// GetFollowers получает подписчиков пользователя с GitHub API
func GetFollowers(username string) ([]string, error) {
	var allFollowers []string
	page := 1

	for {
		url := fmt.Sprintf("%s/users/%s/followers?per_page=%d&page=%d", githubAPI, username, maxFollowersPerPage, page)

		log.Printf("Requesting followers for %s from GitHub API (page %d)", username, page)
		resp, err := makeGitHubAPIRequest(url)
		if err != nil {
			return nil, err
		}

		// Закрываем тело ответа после чтения данных
		defer func(body io.ReadCloser) {
			if err := body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}(resp.Body)

		// Обрабатываем результат запроса
		var followers []struct {
			Login string `json:"login"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
			log.Printf("Error decoding followers from GitHub for %s: %v", username, err)
			return nil, err
		}

		// Добавляем новых подписчиков в общий список
		for _, follower := range followers {
			allFollowers = append(allFollowers, follower.Login)
		}

		// Если количество подписчиков меньше максимального на странице, прекращаем запросы
		if len(followers) < maxFollowersPerPage {
			break
		}

		page++
	}

	log.Printf("Retrieved %d followers for %s from GitHub", len(allFollowers), username)
	return allFollowers, nil
}

// makeGitHubAPIRequest выполняет HTTP-запрос к GitHub API и обрабатывает возможные ошибки
func makeGitHubAPIRequest(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Установка таймаута на запрос
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating GitHub API request: %v", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if githubAPIKey != "" {
		req.Header.Set("Authorization", "token "+githubAPIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making GitHub API request to %s: %v", url, err)
		return nil, err
	}

	// Проверяем на ошибки статуса
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("GitHub API error for %s: %s", url, string(body))
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	log.Printf("GitHub API request to %s succeeded", url)
	return resp, nil
}
