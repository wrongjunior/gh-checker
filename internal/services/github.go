package services

import (
	"encoding/json"
	"fmt"
	"gh-checker/internal/lib/logger"
	"io"
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
	logger.Info("GitHub API key set")
}

// GetFollowers получает подписчиков пользователя с GitHub API
func GetFollowers(username string) ([]string, error) {
	var allFollowers []string
	page := 1

	for {
		url := fmt.Sprintf("%s/users/%s/followers?per_page=%d&page=%d", githubAPI, username, maxFollowersPerPage, page)

		logger.Info(fmt.Sprintf("Requesting followers for %s from GitHub API (page %d)", username, page))
		resp, err := makeGitHubAPIRequest(url)
		if err != nil {
			return nil, err
		}

		// Закрываем тело ответа после чтения данных
		defer func(body io.ReadCloser) {
			if err := body.Close(); err != nil {
				logger.Error("Error closing response body", err)
			}
		}(resp.Body)

		// Обрабатываем результат запроса
		var followers []struct {
			Login string `json:"login"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
			logger.Error(fmt.Sprintf("Error decoding followers from GitHub for %s", username), err)
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

	logger.Info(fmt.Sprintf("Retrieved %d followers for %s from GitHub", len(allFollowers), username))
	return allFollowers, nil
}

// makeGitHubAPIRequest выполняет HTTP-запрос к GitHub API и обрабатывает возможные ошибки
func makeGitHubAPIRequest(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Установка таймаута на запрос
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Error creating GitHub API request", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if githubAPIKey != "" {
		req.Header.Set("Authorization", "token "+githubAPIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Error making GitHub API request to %s", url), err)
		return nil, err
	}

	// Проверяем на ошибки статуса
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Error(fmt.Sprintf("GitHub API error for %s", url), fmt.Errorf("%s", string(body)))
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	logger.Info(fmt.Sprintf("GitHub API request to %s succeeded", url))
	return resp, nil
}
