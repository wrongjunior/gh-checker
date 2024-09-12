package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const githubAPI = "https://api.github.com"

var githubAPIKey string

func SetGitHubAPIKey(apiKey string) {
	githubAPIKey = apiKey
}

func GetFollowers(username string) ([]string, error) {
	url := fmt.Sprintf("%s/users/%s/followers", githubAPI, username)

	resp, err := makeGitHubAPIRequest(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	var followers []struct {
		Login string `json:"login"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
		return nil, err
	}

	var followerNames []string
	for _, follower := range followers {
		followerNames = append(followerNames, follower.Login)
	}

	return followerNames, nil
}

func makeGitHubAPIRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if githubAPIKey != "" {
		req.Header.Set("Authorization", "token "+githubAPIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Проверяем, что GitHub API вернул успешный ответ
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	return resp, nil
}
