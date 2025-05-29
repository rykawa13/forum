package entity

import "errors"

var ErrTopicNotFound = errors.New("topic not found")

type Topic struct {
	ID       string
	Title    string
	Content  string
	UserID   int64
	Username string
}

type Post struct {
	ID       string
	TopicID  string
	Content  string
	UserID   int64
	Username string
}
