package entity

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"-"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
	IsBlocked bool      `db:"is_blocked" json:"is_blocked"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserCreate struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserUpdate struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=6"`
}

type UpdateRoleInput struct {
	IsAdmin bool `json:"is_admin"`
}

type UserListResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

type UserStatus struct {
	IsBlocked bool `json:"is_blocked"`
}

type SessionListResponse struct {
	Sessions []SessionInfo `json:"sessions"`
	Total    int           `json:"total"`
}
