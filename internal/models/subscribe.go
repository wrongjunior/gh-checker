package models

type SubscribeRequest struct {
	Follower string `json:"follower"`
	Followed string `json:"followed"`
}

type SubscribeResponse struct {
	IsFollowing bool   `json:"isFollowing"`
	Error       string `json:"error,omitempty"`
}
