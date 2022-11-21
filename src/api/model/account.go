package model

import (
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	Username       string    `json:"username"` // todo validate email only
	HashedPassword string    `json:"-"`
	GoogleID       string    `json:"google_id"`
	CreateAt       time.Time `json:"create_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewAccount(ctx context.Context, username, password string) (*Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Err(ctx, err)
		return nil, err
	}

	account := &Account{
		Username:       username,
		HashedPassword: string(hashedPassword),
	}

	return account, nil
}

// IsCorrectPassword checks if the provided password is matched with acc's password
func (a *Account) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.HashedPassword), []byte(password))
	return err == nil
}
