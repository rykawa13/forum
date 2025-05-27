package repository

import (
	"auth-service/internal/entity"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SessionPostgres struct {
	db *pgxpool.Pool
}

func NewSessionPostgres(db *pgxpool.Pool) *SessionPostgres {
	return &SessionPostgres{db: db}
}

func (r *SessionPostgres) Create(ctx context.Context, session entity.Session) (int, error) {
	var id int
	query := `
		INSERT INTO sessions (
			user_id, access_token, refresh_token, 
			access_expires, expires_at, user_agent, 
			ip, created_at, updated_at, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), true) 
		RETURNING id`
	err := r.db.QueryRow(ctx, query,
		session.UserID,
		session.AccessToken,
		session.RefreshToken,
		session.AccessExpires,
		session.ExpiresAt,
		session.UserAgent,
		session.IP,
	).Scan(&id)
	return id, err
}

func (r *SessionPostgres) GetByRefreshToken(ctx context.Context, refreshToken string) (entity.Session, error) {
	var session entity.Session
	query := `
		SELECT id, user_id, access_token, refresh_token, 
			   access_expires, expires_at, user_agent, ip, 
			   created_at, updated_at, is_active
		FROM sessions 
		WHERE refresh_token = $1 AND is_active = true`
	err := r.db.QueryRow(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.AccessToken,
		&session.RefreshToken,
		&session.AccessExpires,
		&session.ExpiresAt,
		&session.UserAgent,
		&session.IP,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.IsActive,
	)
	return session, err
}

func (r *SessionPostgres) GetByAccessToken(ctx context.Context, accessToken string) (entity.Session, error) {
	var session entity.Session
	query := `
		SELECT id, user_id, access_token, refresh_token, 
			   access_expires, expires_at, user_agent, ip, 
			   created_at, updated_at, is_active
		FROM sessions 
		WHERE access_token = $1 AND is_active = true`
	err := r.db.QueryRow(ctx, query, accessToken).Scan(
		&session.ID,
		&session.UserID,
		&session.AccessToken,
		&session.RefreshToken,
		&session.AccessExpires,
		&session.ExpiresAt,
		&session.UserAgent,
		&session.IP,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.IsActive,
	)
	return session, err
}

func (r *SessionPostgres) UpdateTokens(ctx context.Context, id int, accessToken, refreshToken string, accessExpires, refreshExpires time.Time) error {
	query := `
		UPDATE sessions 
		SET access_token = $1, refresh_token = $2, 
			access_expires = $3, expires_at = $4, 
			updated_at = NOW()
		WHERE id = $5 AND is_active = true`
	_, err := r.db.Exec(ctx, query, accessToken, refreshToken, accessExpires, refreshExpires, id)
	return err
}

func (r *SessionPostgres) DeactivateSession(ctx context.Context, id int) error {
	query := `
		UPDATE sessions 
		SET is_active = false, updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *SessionPostgres) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *SessionPostgres) DeleteAllUserSessions(ctx context.Context, userID int) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *SessionPostgres) GetUserSessions(ctx context.Context, userID int) ([]entity.Session, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, 
			   access_expires, expires_at, user_agent, ip, 
			   created_at, updated_at, is_active
		FROM sessions 
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []entity.Session
	for rows.Next() {
		var session entity.Session
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.AccessToken,
			&session.RefreshToken,
			&session.AccessExpires,
			&session.ExpiresAt,
			&session.UserAgent,
			&session.IP,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.IsActive,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
