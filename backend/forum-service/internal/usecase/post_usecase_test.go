package usecase

import (
	"context"
	"forum-service/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, input entity.CreatePostInput) (int64, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepository) GetAll(ctx context.Context, offset, limit int) ([]*entity.Post, int, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*entity.Post), args.Int(1), args.Error(2)
}

func (m *MockPostRepository) Update(ctx context.Context, id int64, input entity.UpdatePostInput) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepository) GetReplies(ctx context.Context, postID int64) ([]*entity.Reply, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]*entity.Reply), args.Error(1)
}

func (m *MockPostRepository) CreateReply(ctx context.Context, postID int64, input entity.CreateReplyInput) (*entity.Reply, error) {
	args := m.Called(ctx, postID, input)
	return args.Get(0).(*entity.Reply), args.Error(1)
}

func (m *MockPostRepository) DeleteReply(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestPostUseCase_Create(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockUserRepo := new(MockUserRepository)
	useCase := NewPostUseCase(mockRepo, mockUserRepo)
	ctx := context.Background()

	input := entity.CreatePostInput{
		Title:    "Test Post",
		Content:  "Test Content",
		AuthorID: 1,
	}

	expectedID := int64(1)
	mockRepo.On("Create", ctx, input).Return(expectedID, nil)

	id, err := useCase.Create(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	mockRepo.AssertExpectations(t)
}

func TestPostUseCase_GetByID(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockUserRepo := new(MockUserRepository)
	useCase := NewPostUseCase(mockRepo, mockUserRepo)
	ctx := context.Background()

	expectedPost := &entity.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(expectedPost, nil)

	post, err := useCase.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
	mockRepo.AssertExpectations(t)
}
