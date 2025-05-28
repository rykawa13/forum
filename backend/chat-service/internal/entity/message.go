package entity

import "time"

// Message представляет сообщение в чате
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMessage создает новое сообщение
func NewMessage(content string, userID int64, username string) *Message {
	now := time.Now()
	return &Message{
		Content:   content,
		UserID:    userID,
		Username:  username,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type MessageCreate struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}
