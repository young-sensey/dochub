package entities

import "time"

type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
	} `json:"user"`
}
