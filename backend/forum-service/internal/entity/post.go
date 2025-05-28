package entity

import "time"

// Post represents a forum post
type Post struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int64     `json:"author_id" db:"author_id"`
	Author    *User     `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Reply represents a reply to a post
type Reply struct {
	ID        int64     `json:"id" db:"id"`
	PostID    int64     `json:"post_id" db:"post_id"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int64     `json:"author_id" db:"author_id"`
	Author    *User     `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePostInput represents input data for creating a new post
type CreatePostInput struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Content  string `json:"content" binding:"required,min=1"`
	AuthorID int64  `json:"author_id" binding:"required"`
}

// UpdatePostInput represents input data for updating a post
type UpdatePostInput struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=255"`
	Content string `json:"content" binding:"omitempty,min=1"`
}

// CreateReplyInput represents input data for creating a new reply
type CreateReplyInput struct {
	Content  string `json:"content" binding:"required,min=1"`
	AuthorID int64  `json:"author_id" binding:"required"`
}

type CreatePostDTO struct {
	Title    string `json:"title" validate:"required,min=5,max=200"`
	Content  string `json:"content" validate:"required,min=10"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type UpdatePostDTO struct {
	Title    *string `json:"title,omitempty" validate:"omitempty,min=5,max=200"`
	Content  *string `json:"content,omitempty" validate:"omitempty,min=10"`
	IsLocked *bool   `json:"is_locked,omitempty"`
}

type CreateReplyDTO struct {
	PostID   int64  `json:"post_id"`
	Content  string `json:"content" validate:"required,min=1"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type ForumStats struct {
	TotalPosts    int `json:"total_posts"`
	TotalReplies  int `json:"total_replies"`
	ActiveThreads int `json:"active_threads"`
	UsersPosted   int `json:"users_posted"`
}
