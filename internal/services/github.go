package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const githubAPI = "https://api.github.com"

var githubAPIKey string

func SetGitHubAPIKey(apiKey string) {
	githubAPIKey = apiKey
	log.Println("GitHub API key set")
}

func GetFollowers(username string) ([]string, error) {
	url := fmt.Sprintf("%s/users/%s/followers", githubAPI, username)

	log.Printf("Requesting followers for %s from GitHub API", username)
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
		log.Printf("Error decoding followers from GitHub for %s: %v", username, err)
		return nil, err
	}

	var followerNames []string
	for _, follower := range followers {
		followerNames = append(followerNames, follower.Login)
	}

	log.Printf("Retrieved %d followers for %s from GitHub", len(followerNames), username)
	return followerNames, nil
}

func makeGitHubAPIRequest(url string) (*http.Response, error) {
	client := &http.Client{}
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("GitHub API error for %s: %s", url, string(body))
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	log.Printf("GitHub API request to %s succeeded", url)
	return resp, nil
}
