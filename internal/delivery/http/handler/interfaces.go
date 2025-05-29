package handler

import "backend/forum-service/internal/entity"

type ForumUseCase interface {
	CreateTopic(topic *entity.Topic) error
	GetTopic(id string) (*entity.Topic, error)
	ListTopics(page, limit int) ([]*entity.Topic, error)
	CreatePost(post *entity.Post) error
}
