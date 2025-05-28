package repository

import (
	"context"
	"forum-service/internal/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, username
		FROM users
		WHERE id = $1
	`

	user := &entity.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
