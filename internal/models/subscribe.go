package models

type SubscribeRequest struct {
	Follower string `json:"follower"`
	Followed string `json:"followed"`
}

type SubscribeResponse struct {
	IsFollowing bool   `json:"isFollowing"`
	Error       string `json:"error,omitempty"`
}

type StarCheckRequest struct {
	Username   string `json:"username"`
	Repository string `json:"repository"`
}

type StarCheckResponse struct {
	HasStar bool   `json:"hasStar"`
	Error   string `json:"error,omitempty"`
}
