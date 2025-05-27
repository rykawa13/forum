package repository

import (
	"auth-service/internal/entity"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StatsRepo struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepo {
	return &StatsRepo{db: db}
}

func (r *StatsRepo) GetStats(ctx context.Context) (*entity.ForumStats, error) {
	stats := &entity.ForumStats{}

	// Получаем количество пользователей
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		return nil, err
	}

	// Для демонстрации, пока нет реальных данных
	stats.TotalTopics = 0
	stats.TotalPosts = 0

	return stats, nil
}
