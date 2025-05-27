package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return NewUserPostgres(db)
}

// NewSessionRepository creates a new instance of SessionRepository
func NewSessionRepository(db *pgxpool.Pool) ISessionRepository {
	return NewSessionPostgres(db)
}
