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
	Username   string `json:"username"`   // Пользователь, который ставит звезду
	Repository string `json:"repository"` // Репозиторий, на который ставится звезда
}

type StarCheckResponse struct {
	HasStar bool   `json:"hasStar"` // Флаг: есть ли звезда на репозитории
	Error   string `json:"error,omitempty"`
}
