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

func CheckIfFollowing(follower, followed string) (bool, error) {
	url := fmt.Sprintf("%s/users/%s/following/%s", githubAPI, follower, followed)
	return makeGitHubRequest(url)
}

func GetFollowers(username string) ([]string, error) {
	url := fmt.Sprintf("%s/users/%s/followers", githubAPI, username)

	resp, err := makeGitHubAPIRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

func makeGitHubRequest(url string) (bool, error) {
	resp, err := makeGitHubAPIRequest(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("GitHub API error: %s", string(body))
	}
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

	return client.Do(req)
}
