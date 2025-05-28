package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitDB initializes the database schema
func InitDB(ctx context.Context, pool *pgxpool.Pool) error {
	// Create posts table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author_id BIGINT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
		
		CREATE TABLE IF NOT EXISTS replies (
			id BIGSERIAL PRIMARY KEY,
			post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			author_id BIGINT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_replies_post_id ON replies(post_id);
		CREATE INDEX IF NOT EXISTS idx_replies_author_id ON replies(author_id);
	`)
	return err
}
