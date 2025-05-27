package repository

import (
	"auth-service/internal/entity"
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)`

	now := time.Now()

	// Проверяем хэш пароля перед вставкой
	if len(user.Password) < 20 {
		log.Printf("Warning: Password hash seems too short: %d chars", len(user.Password))
	} else {
		log.Printf("Inserting user with password hash length: %d", len(user.Password))
	}

	_, err := r.db.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.IsAdmin,
		now,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return errors.New("user with this email already exists")
		}
		log.Printf("Error creating user: %v", err)
		return err
	}

	log.Printf("Successfully created user")
	return nil
}

func (r *UserPostgres) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user := &entity.User{}
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserPostgres) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		log.Printf("Database error when getting user by email: %v", err)
		return nil, err
	}

	// Проверяем, что получили хэш пароля
	if user.Password == "" {
		log.Printf("Warning: Empty password hash for user %s", email)
	} else {
		log.Printf("Successfully retrieved password hash for user %s", email)
	}

	return user, nil
}

func (r *UserPostgres) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET username = $1, 
			email = $2, 
			password_hash = $3, 
			is_admin = $4, 
			updated_at = $5
		WHERE id = $6`

	result, err := r.db.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.IsAdmin,
		time.Now(),
		user.ID,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return errors.New("email already taken")
		}
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserPostgres) GetAll(ctx context.Context) ([]*entity.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserPostgres) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
