package model

import "time"

type Account struct {
	Username   string    `json:"username"` // todo validate email only
	Password   string    `json:"password"`
	GoogleID   string    `json:"google_id"`
	LastLogout time.Time `json:"last_logout"`
	CreateAt   time.Time `json:"create_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}