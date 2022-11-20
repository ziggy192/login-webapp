package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	Username       string    `json:"username"` // todo validate email only
	HashedPassword string    `json:"-"`
	GoogleID       string    `json:"google_id"`
	LastLogout     time.Time `json:"last_logout"`
	CreateAt       time.Time `json:"create_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// IsCorrectPassword checks if the provided password is matched with acc's password
func (a *Account) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.HashedPassword), []byte(password))
	return err == nil
}
