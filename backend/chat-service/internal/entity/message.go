package entity

import "time"

type Message struct {
	ID        int       `db:"id" json:"id"`
	Content   string    `db:"content" json:"content"`
	AuthorID  int       `db:"author_id" json:"author_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type MessageCreate struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}
