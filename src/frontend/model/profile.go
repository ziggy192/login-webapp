package model

import "time"

type Profile struct {
	FullName  string    `json:"full_name" schema:"full_name"`
	Phone     string    `json:"phone" schema:"phone"`
	Email     string    `json:"email" schema:"email"`
	AccountID string    `json:"account_id" schema:"-"`
	CreatedAt time.Time `json:"created_at" schema:"-"`
	UpdatedAt time.Time `json:"updated_at" schema:"-"`
}
