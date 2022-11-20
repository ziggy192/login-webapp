package model

type Account struct {
	Username string `json:"username" schema:"username"` // todo validate email only
	Password string `json:"password" schema:"password"`
}
