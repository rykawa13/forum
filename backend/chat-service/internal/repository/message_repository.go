package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/chat-service/internal/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

// MessageRepository интерфейс для работы с сообщениями
type MessageRepository interface {
	Create(ctx context.Context, message *entity.Message) error
	GetHistory(ctx context.Context, limit int32, beforeID int64) ([]*entity.Message, error)
	DeleteOldMessages(ctx context.Context, before time.Time) (int32, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
}

type Row interface {
	Scan(dest ...interface{}) error
}

type messageRepository struct {
	pool *pgxpool.Pool
}

// NewMessageRepository создает новый репозиторий сообщений
func NewMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) Create(ctx context.Context, message *entity.Message) error {
	// Проверяем состояние соединения с БД
	if err := r.pool.Ping(ctx); err != nil {
		log.Printf("Database connection error: %v", err)
		return fmt.Errorf("database connection error: %v", err)
	}

	query := `
        INSERT INTO messages (content, user_id, username, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id::text`

	log.Printf("Executing query: %s with values: content=%s, user_id=%d, username=%s",
		query, message.Content, message.UserID, message.Username)

	// Используем контекст для выполнения запроса
	row := r.pool.QueryRow(ctx, query,
		message.Content,
		message.UserID,
		message.Username,
		message.CreatedAt,
		message.UpdatedAt,
	)

	err := row.Scan(&message.ID)
	if err != nil {
		if ctx.Err() != nil {
			// Если ошибка связана с контекстом
			log.Printf("Context error during query execution: %v", ctx.Err())
			return ctx.Err()
		}
		log.Printf("Error executing query: %v", err)
		return fmt.Errorf("database error: %v", err)
	}

	log.Printf("Successfully created message with ID: %s", message.ID)
	return nil
}

func (r *messageRepository) GetHistory(ctx context.Context, limit int32, beforeID int64) ([]*entity.Message, error) {
	query := `
        SELECT id::text, content, user_id, username, created_at, updated_at
        FROM messages
        WHERE ($2 = 0 OR id < $2)
        ORDER BY created_at ASC
        LIMIT $1`

	log.Printf("Executing query: %s with values: limit=%d, beforeID=%d", query, limit, beforeID)

	rows, err := r.pool.Query(ctx, query, limit, beforeID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []*entity.Message
	for rows.Next() {
		msg := &entity.Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.UserID,
			&msg.Username,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("failed to iterate messages: %v", err)
	}

	log.Printf("Successfully retrieved %d messages", len(messages))
	return messages, nil
}

func (r *messageRepository) DeleteOldMessages(ctx context.Context, before time.Time) (int32, error) {
	query := `
        DELETE FROM messages
        WHERE created_at < $1`

	log.Printf("Executing query: %s with value: before=%v", query, before)

	result, err := r.pool.Exec(ctx, query, before)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return 0, err
	}

	rowsAffected := int32(result.RowsAffected())
	log.Printf("Successfully deleted %d old messages", rowsAffected)
	return rowsAffected, nil
}

func (r *messageRepository) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return r.pool.QueryRow(ctx, query, args...)
}
