package usecase

import (
	"auth-service/internal/entity"
	"auth-service/internal/repository"
	"auth-service/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.ISessionRepository
	tokenManager jwt.TokenManager
	config       *Config
}

type Config struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewAuthUseCase(userRepo repository.UserRepository, sessionRepo repository.ISessionRepository, tokenManager jwt.TokenManager, config *Config) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
		config:       config,
	}
}

func (u *AuthUseCase) CreateUser(ctx context.Context, input entity.UserCreate) error {
	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	user := &entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	return u.userRepo.Create(ctx, user)
}

func (u *AuthUseCase) SignIn(ctx context.Context, input *entity.SignInInput, userAgent, ip string) (*entity.Tokens, error) {
	// Получаем пользователя по email
	log.Printf("Attempting to sign in user with email: %s", input.Email)

	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		log.Printf("Error finding user by email %s: %v", input.Email, err)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("User found: %s", user.Email)

	// Сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		log.Printf("Password comparison failed for user %s: %v", input.Email, err)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("Password verified successfully for user: %s", user.Email)

	// Если все хорошо, создаем сессию
	return u.createSession(ctx, user.ID, userAgent, ip)
}

func (u *AuthUseCase) RefreshTokens(ctx context.Context, refreshToken, userAgent, ip string) (*entity.Tokens, error) {
	session, err := u.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	return u.createSession(ctx, session.UserID, userAgent, ip)
}

func (u *AuthUseCase) createSession(ctx context.Context, userID int, userAgent, ip string) (*entity.Tokens, error) {
	// Создаем новые токены
	accessToken, err := u.tokenManager.NewJWT(fmt.Sprintf("%d", userID), u.config.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	session := entity.Session{
		UserID:        userID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpires: time.Now().Add(u.config.AccessTokenTTL),
		ExpiresAt:     time.Now().Add(u.config.RefreshTokenTTL),
		UserAgent:     userAgent,
		IP:            ip,
		IsActive:      true,
	}

	// Деактивируем все предыдущие сессии пользователя
	if err := u.sessionRepo.DeleteAllUserSessions(ctx, userID); err != nil {
		return nil, err
	}

	// Создаем новую сессию
	_, err = u.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	session, err := u.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return u.sessionRepo.DeactivateSession(ctx, session.ID)
}

func (u *AuthUseCase) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *AuthUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return u.userRepo.Update(ctx, user)
}

func (u *AuthUseCase) UpdateUserRole(ctx context.Context, userID int, isAdmin bool) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.IsAdmin = isAdmin
	return u.userRepo.Update(ctx, user)
}

func (u *AuthUseCase) GetUserSessions(ctx context.Context, userID int) ([]entity.SessionInfo, int, error) {
	sessions, err := u.sessionRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	sessionInfos := make([]entity.SessionInfo, len(sessions))
	for i, session := range sessions {
		sessionInfos[i] = entity.SessionInfo{
			ID:        session.ID,
			UserID:    session.UserID,
			UserAgent: session.UserAgent,
			IP:        session.IP,
			CreatedAt: session.CreatedAt,
			IsActive:  session.IsActive,
		}
	}

	return sessionInfos, len(sessionInfos), nil
}

func (u *AuthUseCase) TerminateUserSessions(ctx context.Context, userID int) error {
	return u.sessionRepo.DeleteAllUserSessions(ctx, userID)
}

// ParseToken парсит и проверяет JWT токен
func (u *AuthUseCase) ParseToken(ctx context.Context, token string) (*entity.User, error) {
	// Парсим токен
	userID, err := u.tokenManager.Parse(token)
	if err != nil {
		return nil, err
	}

	// Проверяем сессию
	session, err := u.sessionRepo.GetByAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !session.IsActive || time.Now().After(session.AccessExpires) {
		return nil, errors.New("token expired")
	}

	// Получаем пользователя
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	return u.userRepo.GetByID(ctx, id)
}

// SignUp регистрирует нового пользователя
func (u *AuthUseCase) SignUp(ctx context.Context, input entity.UserCreate) error {
	return u.CreateUser(ctx, input)
}

// GetAll возвращает список всех пользователей
func (u *AuthUseCase) GetAll(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *AuthUseCase) DeleteUser(ctx context.Context, userID int) error {
	// Проверяем существование пользователя
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем, не пытаемся ли удалить последнего админа
	if user.IsAdmin {
		admins, err := u.userRepo.GetAll(ctx)
		if err != nil {
			return err
		}
		adminCount := 0
		for _, u := range admins {
			if u.IsAdmin {
				adminCount++
			}
		}
		if adminCount <= 1 {
			return errors.New("cannot delete the last admin user")
		}
	}

	// Удаляем все сессии пользователя
	if err := u.sessionRepo.DeleteAllUserSessions(ctx, userID); err != nil {
		return err
	}

	// Удаляем пользователя
	return u.userRepo.Delete(ctx, userID)
}
