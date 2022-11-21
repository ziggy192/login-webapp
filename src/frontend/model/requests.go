package model

type LoginRequest struct {
	Username string `json:"username"` // todo validate email only
	Password string `json:"password"`
}
