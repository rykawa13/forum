package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/backend/auth-service/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (int, error)
	GetByID(ctx context.Context, id int) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, user entity.User) (int, error) {
	var id int
	query := `INSERT INTO users (username, email, password_hash, is_admin) 
	          VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.Password, user.IsAdmin).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UserPostgres.Create: %w", err)
	}

	return id, nil
}

func (r *UserPostgres) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	query := `SELECT id, username, email, password_hash, is_admin, created_at, updated_at 
	          FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, fmt.Errorf("UserPostgres.GetByEmail: %w", entity.ErrUserNotFound)
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("UserPostgres.GetByEmail: %w", err)
	}

	return user, nil
}
