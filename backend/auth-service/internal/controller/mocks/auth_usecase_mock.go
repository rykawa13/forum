package mocks

import (
	"context"

	"github.com/forum-backend/auth-service/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, input entity.UserCreate) (entity.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockAuthUseCase) Login(ctx context.Context, input entity.UserLogin) (entity.Tokens, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(entity.Tokens), args.Error(1)
}
