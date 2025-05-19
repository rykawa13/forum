package entity

import "time"

type Session struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	UserAgent    string    `db:"user_agent"`
	IP           string    `db:"ip"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
