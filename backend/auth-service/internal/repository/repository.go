package repository

import (
	"auth-service/internal/entity"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
}

type StatsRepository interface {
	GetStats(ctx context.Context) (*entity.ForumStats, error)
}

type ISessionRepository interface {
	Create(ctx context.Context, session entity.Session) (int, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (entity.Session, error)
	GetByAccessToken(ctx context.Context, accessToken string) (entity.Session, error)
	Delete(ctx context.Context, id int) error
	DeleteAllUserSessions(ctx context.Context, userID int) error
	UpdateTokens(ctx context.Context, id int, accessToken, refreshToken string, accessExpires, refreshExpires time.Time) error
	DeactivateSession(ctx context.Context, id int) error
	GetUserSessions(ctx context.Context, userID int) ([]entity.Session, error)
}

type Repository struct {
	UserRepository
	StatsRepository
	SessionRepository ISessionRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		UserRepository:    NewUserRepository(db),
		StatsRepository:   NewStatsRepository(db),
		SessionRepository: NewSessionRepository(db),
	}
}
