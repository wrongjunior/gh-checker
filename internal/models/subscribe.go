package models

// Модель для запроса подписки
type SubscribeRequest struct {
	Follower string `json:"follower"`
	Followed string `json:"followed"`
}

// Модель для ответа
type SubscribeResponse struct {
	IsFollowing bool   `json:"isFollowing"`
	Error       string `json:"error,omitempty"`
}
