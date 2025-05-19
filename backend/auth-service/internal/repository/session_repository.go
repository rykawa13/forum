package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/forum-backend/auth-service/internal/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, session entity.Session) (int, error)
	GetByToken(ctx context.Context, refreshToken string) (entity.Session, error)
	Delete(ctx context.Context, id int) error
	DeleteExpired(ctx context.Context) error
}

type SessionPostgres struct {
	db *sql.DB
}

func NewSessionPostgres(db *sql.DB) *SessionPostgres {
	return &SessionPostgres{db: db}
}

func (r *SessionPostgres) Create(ctx context.Context, session entity.Session) (int, error) {
	var id int
	query := `INSERT INTO sessions (user_id, refresh_token, user_agent, ip, expires_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IP,
		session.ExpiresAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("SessionPostgres.Create: %w", err)
	}

	return id, nil
}

func (r *SessionPostgres) GetByToken(ctx context.Context, refreshToken string) (entity.Session, error) {
	var session entity.Session
	query := `SELECT id, user_id, refresh_token, user_agent, ip, expires_at, created_at
	          FROM sessions WHERE refresh_token = $1`

	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IP,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Session{}, fmt.Errorf("SessionPostgres.GetByToken: %w", entity.ErrSessionNotFound)
	}
	if err != nil {
		return entity.Session{}, fmt.Errorf("SessionPostgres.GetByToken: %w", err)
	}

	return session, nil
}

func (r *SessionPostgres) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("SessionPostgres.Delete: %w", err)
	}
	return nil
}

func (r *SessionPostgres) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("SessionPostgres.DeleteExpired: %w", err)
	}
	return nil
}
