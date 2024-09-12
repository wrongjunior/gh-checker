package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const githubAPI = "https://api.github.com"

var githubAPIKey string

// SetGitHubAPIKey задает GitHub API ключ для аутентификации запросов
func SetGitHubAPIKey(apiKey string) {
	githubAPIKey = apiKey
}

// CheckIfFollowing проверяет, подписан ли один пользователь на другого
func CheckIfFollowing(follower, followed string) (bool, error) {
	url := fmt.Sprintf("%s/users/%s/following/%s", githubAPI, follower, followed)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Добавляем API ключ в заголовок, если он задан
	if githubAPIKey != "" {
		req.Header.Set("Authorization", "token "+githubAPIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("GitHub API error: %s", string(body))
	}
}
