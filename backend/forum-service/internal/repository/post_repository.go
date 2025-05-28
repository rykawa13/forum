package repository

import (
	"context"
	"forum-service/internal/entity"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	GetAll(ctx context.Context, offset, limit int) ([]*entity.Post, int, error)
	Create(ctx context.Context, input entity.CreatePostInput) (int64, error)
	Update(ctx context.Context, id int64, input entity.UpdatePostInput) error
	Delete(ctx context.Context, id int64) error
	GetReplies(ctx context.Context, postID int64) ([]*entity.Reply, error)
	CreateReply(ctx context.Context, postID int64, input entity.CreateReplyInput) (*entity.Reply, error)
	DeleteReply(ctx context.Context, id int64) error
}

type postRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) PostRepository {
	return &postRepository{
		pool: pool,
	}
}

func (r *postRepository) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	post := &entity.Post{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepository) GetAll(ctx context.Context, offset, limit int) ([]*entity.Post, int, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		post := &entity.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	var total int
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM posts").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Create(ctx context.Context, input entity.CreatePostInput) (int64, error) {
	query := `
		INSERT INTO posts (title, content, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id
	`

	now := time.Now()
	var id int64
	err := r.pool.QueryRow(ctx, query,
		input.Title,
		input.Content,
		input.AuthorID,
		now,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *postRepository) Update(ctx context.Context, id int64, input entity.UpdatePostInput) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = $3
		WHERE id = $4
	`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query,
		input.Title,
		input.Content,
		now,
		id,
	)
	return err
}

func (r *postRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *postRepository) GetReplies(ctx context.Context, postID int64) ([]*entity.Reply, error) {
	query := `
		SELECT id, post_id, content, author_id, created_at, updated_at
		FROM replies
		WHERE post_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*entity.Reply
	for rows.Next() {
		reply := &entity.Reply{}
		err := rows.Scan(
			&reply.ID,
			&reply.PostID,
			&reply.Content,
			&reply.AuthorID,
			&reply.CreatedAt,
			&reply.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}

	return replies, nil
}

func (r *postRepository) CreateReply(ctx context.Context, postID int64, input entity.CreateReplyInput) (*entity.Reply, error) {
	query := `
		INSERT INTO replies (post_id, content, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	reply := &entity.Reply{
		PostID:   postID,
		Content:  input.Content,
		AuthorID: input.AuthorID,
	}

	err := r.pool.QueryRow(ctx, query,
		postID,
		input.Content,
		input.AuthorID,
		now,
	).Scan(&reply.ID, &reply.CreatedAt, &reply.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (r *postRepository) DeleteReply(ctx context.Context, id int64) error {
	query := `DELETE FROM replies WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
