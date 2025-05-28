package usecase

import (
	"context"
	"errors"
	"forum-service/internal/entity"
	"forum-service/internal/repository"
)

type PostUseCase interface {
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	GetAll(ctx context.Context, offset, limit int) ([]*entity.Post, int, error)
	Create(ctx context.Context, input entity.CreatePostInput) (*entity.Post, error)
	Update(ctx context.Context, id int64, input entity.UpdatePostInput) error
	Delete(ctx context.Context, id int64) error
	GetReplies(ctx context.Context, postID int64) ([]*entity.Reply, error)
	CreateReply(ctx context.Context, postID int64, input entity.CreateReplyInput) (*entity.Reply, error)
	DeleteReply(ctx context.Context, id int64) error
}

type postUseCase struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostUseCase(postRepo repository.PostRepository, userRepo repository.UserRepository) PostUseCase {
	return &postUseCase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (uc *postUseCase) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	author, err := uc.userRepo.GetByID(ctx, post.AuthorID)
	if err != nil {
		// Если не удалось получить информацию об авторе, продолжаем без неё
		return post, nil
	}

	post.Author = author
	return post, nil
}

func (uc *postUseCase) GetAll(ctx context.Context, offset, limit int) ([]*entity.Post, int, error) {
	posts, total, err := uc.postRepo.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Получаем информацию об авторах для всех постов
	for _, post := range posts {
		author, err := uc.userRepo.GetByID(ctx, post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	return posts, total, nil
}

func (uc *postUseCase) Create(ctx context.Context, input entity.CreatePostInput) (*entity.Post, error) {
	if input.Title == "" || input.Content == "" {
		return nil, errors.New("title and content are required")
	}

	// Create post and get its ID
	id, err := uc.postRepo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	// Get the created post
	return uc.GetByID(ctx, id)
}

func (uc *postUseCase) Update(ctx context.Context, id int64, input entity.UpdatePostInput) error {
	_, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.postRepo.Update(ctx, id, input)
}

func (uc *postUseCase) Delete(ctx context.Context, id int64) error {
	return uc.postRepo.Delete(ctx, id)
}

func (uc *postUseCase) GetReplies(ctx context.Context, postID int64) ([]*entity.Reply, error) {
	replies, err := uc.postRepo.GetReplies(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Получаем информацию об авторах для всех ответов
	for _, reply := range replies {
		author, err := uc.userRepo.GetByID(ctx, reply.AuthorID)
		if err == nil {
			reply.Author = author
		}
	}

	return replies, nil
}

func (uc *postUseCase) CreateReply(ctx context.Context, postID int64, input entity.CreateReplyInput) (*entity.Reply, error) {
	if input.Content == "" {
		return nil, errors.New("reply content is required")
	}

	// Check if post exists
	_, err := uc.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	reply, err := uc.postRepo.CreateReply(ctx, postID, input)
	if err != nil {
		return nil, err
	}

	author, err := uc.userRepo.GetByID(ctx, reply.AuthorID)
	if err == nil {
		reply.Author = author
	}

	return reply, nil
}

func (uc *postUseCase) DeleteReply(ctx context.Context, id int64) error {
	return uc.postRepo.DeleteReply(ctx, id)
}
