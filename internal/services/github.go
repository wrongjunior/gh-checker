package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const githubAPI = "https://api.github.com"

var githubAPIKey string

func SetGitHubAPIKey(apiKey string) {
	githubAPIKey = apiKey
}

func CheckIfFollowing(follower, followed string) (bool, error) {
	url := fmt.Sprintf("%s/users/%s/following/%s", githubAPI, follower, followed)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if githubAPIKey != "" {
		req.Header.Set("Authorization", "token "+githubAPIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("GitHub API error: %s", string(body))
	}
}

func GetFollowers(username string) ([]string, error) {
	url := fmt.Sprintf("%s/users/%s/followers", githubAPI, username)

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

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
