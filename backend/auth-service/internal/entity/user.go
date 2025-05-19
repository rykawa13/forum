package entity

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"-"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserCreate struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
