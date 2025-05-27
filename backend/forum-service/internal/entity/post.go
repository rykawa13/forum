package entity

import "time"

type Post struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	AuthorID  int       `db:"author_id" json:"author_id"`
	Category  string    `db:"category" json:"category"`
	IsVisible bool      `db:"is_visible" json:"is_visible"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type PostCreate struct {
	Title    string `json:"title" validate:"required,min=3,max=100"`
	Content  string `json:"content" validate:"required,min=10"`
	Category string `json:"category" validate:"required"`
}
