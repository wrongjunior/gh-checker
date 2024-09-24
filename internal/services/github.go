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

	logger.Info("Starting to fetch followers for user " + username)

	for {
		url := fmt.Sprintf("%s/users/%s/followers?per_page=%d&page=%d", githubAPI, username, maxFollowersPerPage, page)

		logger.Info(fmt.Sprintf("Requesting followers for %s from GitHub API (page %d)", username, page))
		resp, err := makeGitHubAPIRequestWithRetries(url)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get followers for %s (page %d)", username, page), err)
			return nil, err
		}

		// Обрабатываем результат запроса
		var followers []struct {
			Login string `json:"login"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
			logger.Error(fmt.Sprintf("Error decoding followers from GitHub for %s", username), err)
			resp.Body.Close() // Закрытие тела ответа при ошибке
			return nil, err
		}

		resp.Body.Close() // Закрытие тела после успешного получения данных

		// Добавляем новых подписчиков в общий список
		for _, follower := range followers {
			allFollowers = append(allFollowers, follower.Login)
		}

		// Логируем, сколько подписчиков было обработано на текущей странице
		logger.Info(fmt.Sprintf("Processed %d followers for %s from GitHub (page %d)", len(followers), username, page))

		// Если количество подписчиков меньше максимального на странице, значит больше страниц нет
		if len(followers) < maxFollowersPerPage {
			break
		}

		// Увеличиваем номер страницы
		page++
	}

	logger.Info(fmt.Sprintf("Retrieved %d followers for %s from GitHub", len(allFollowers), username))
	return allFollowers, nil
}

// makeGitHubAPIRequest выполняет HTTP-запрос к GitHub API и обрабатывает возможные ошибки с повторными попытками
func makeGitHubAPIRequestWithRetries(url string) (*http.Response, error) {
	var resp *http.Response
	var err error
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.Info(fmt.Sprintf("Attempt %d to make GitHub API request to %s", attempt, url))
		resp, err = makeGitHubAPIRequest(url)
		if err == nil {
			return resp, nil
		}
		logger.Error(fmt.Sprintf("Error making GitHub API request to %s (attempt %d)", url, attempt), err)

		// Если последняя попытка, возвращаем ошибку
		if attempt == maxAttempts {
			return nil, err
		}

		// Ждем перед повторной попыткой
		time.Sleep(2 * time.Second)
	}

	return nil, err
}

// makeGitHubAPIRequest выполняет HTTP-запрос к GitHub API и обрабатывает возможные ошибки
func makeGitHubAPIRequest(url string) (*http.Response, error) {
	logger.Info("Making GitHub API request to " + url)

	client := &http.Client{
		Timeout: 15 * time.Second, // Увеличенный таймаут на запрос
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
		if closeErr := resp.Body.Close(); closeErr != nil { // Обработка ошибки при закрытии
			logger.Error("Error closing response body after API error", closeErr)
		}
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	logger.Info(fmt.Sprintf("GitHub API request to %s succeeded", url))
	return resp, nil
}

// CheckStar проверяет, поставил ли пользователь звезду на репозиторий
func CheckStar(username, repository string) (bool, error) {
	page := 1

	for {
		url := fmt.Sprintf("%s/repos/%s/stargazers?per_page=%d&page=%d", githubAPI, repository, maxFollowersPerPage, page)

		logger.Info(fmt.Sprintf("Checking if user %s starred repository %s (page %d)", username, repository, page))
		resp, err := makeGitHubAPIRequestWithRetries(url)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to check star for user %s on repository %s", username, repository), err)
			return false, err
		}

		var stargazers []struct {
			Login string `json:"login"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&stargazers); err != nil {
			logger.Error(fmt.Sprintf("Error decoding stargazers from GitHub for %s", repository), err)
			resp.Body.Close() // Закрываем тело ответа при ошибке
			return false, err
		}

		resp.Body.Close() // Закрываем тело после успешного получения данных

		// Проверяем, есть ли пользователь среди тех, кто поставил звезду
		for _, stargazer := range stargazers {
			if stargazer.Login == username {
				logger.Info(fmt.Sprintf("User %s has starred repository %s", username, repository))
				return true, nil
			}
		}

		// Если количество "звезд" меньше максимального на странице, значит больше страниц нет
		if len(stargazers) < maxFollowersPerPage {
			break
		}

		// Переход к следующей странице
		page++
	}

	logger.Info(fmt.Sprintf("User %s has not starred repository %s", username, repository))
	return false, nil
}
