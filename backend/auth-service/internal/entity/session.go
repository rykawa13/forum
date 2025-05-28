package entity

import "time"

type Session struct {
	ID            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
	AccessToken   string    `db:"access_token" json:"access_token"`
	RefreshToken  string    `db:"refresh_token" json:"refresh_token"`
	AccessExpires time.Time `db:"access_expires" json:"access_expires"`
	ExpiresAt     time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	UserAgent     string    `db:"user_agent" json:"user_agent"`
	IP            string    `db:"ip" json:"ip"`
	IsActive      bool      `db:"is_active" json:"is_active"`
}

type SessionInfo struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}
