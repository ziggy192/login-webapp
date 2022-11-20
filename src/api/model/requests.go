package model

import (
	"fmt"
	"net/mail"
)

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *SignupRequest) Validate() error {
	if len(l.Username) == 0 || len(l.Password) == 0 {
		return fmt.Errorf("empty username or password")
	}
	_, err := mail.ParseAddress(l.Username)
	if err != nil {
		return fmt.Errorf("invalid email")
	}
	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
