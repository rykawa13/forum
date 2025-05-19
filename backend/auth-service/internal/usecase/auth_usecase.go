package usecase

import (
	"context"
	"time"

	"github.com/forum-backend/auth-service/internal/entity"
	"github.com/forum-backend/auth-service/internal/repository"
	"github.com/forum-backend/pkg/jwt"
	"github.com/forum-backend/pkg/password"
)

type AuthUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthUseCase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, input entity.UserCreate) (entity.User, error) {
	// Проверка существования пользователя
	_, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return entity.User{}, entity.ErrUserAlreadyExists
	}

	// Хеширование пароля
	hashedPass, err := password.HashPassword(input.Password)
	if err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPass,
		IsAdmin:  false,
	}

	// Создание пользователя
	userID, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	user.ID = userID
	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input entity.UserLogin) (entity.Tokens, error) {
	// Получение пользователя
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return entity.Tokens{}, err
	}

	// Проверка пароля
	if !password.CheckPasswordHash(input.Password, user.Password) {
		return entity.Tokens{}, entity.ErrInvalidCredentials
	}

	// Генерация токенов
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.IsAdmin)
	if err != nil {
		return entity.Tokens{}, err
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return entity.Tokens{}, err
	}

	// Сохранение сессии
	session := entity.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}

	if _, err := uc.sessionRepo.Create(ctx, session); err != nil {
		return entity.Tokens{}, err
	}

	return entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUseCase) RefreshTokens(ctx context.Context, refreshToken string) (entity.Tokens, error) {
	// Получаем сессию по refresh токену
	session, err := uc.sessionRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return entity.Tokens{}, err
	}

	// Проверяем не истек ли токен
	if time.Now().After(session.ExpiresAt) {
		return entity.Tokens{}, entity.ErrRefreshTokenExpired
	}

	// Получаем пользователя
	user, err := uc.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return entity.Tokens{}, err
	}

	// Генерируем новые токены
	newAccessToken, err := jwt.GenerateAccessToken(user.ID, user.IsAdmin)
	if err != nil {
		return entity.Tokens{}, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return entity.Tokens{}, err
	}

	// Обновляем сессию
	err = uc.sessionRepo.Delete(ctx, session.ID)
	if err != nil {
		return entity.Tokens{}, err
	}

	newSession := entity.Session{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}

	if _, err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return entity.Tokens{}, err
	}

	return entity.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
