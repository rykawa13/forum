# Forum Service

## Description
Forum service provides REST and gRPC APIs for managing forum posts and replies.

## Prerequisites
- Go 1.19 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Configuration
Configuration is done through environment variables in `.env` file:
- `DB_URL` - PostgreSQL connection string
- `HTTP_PORT` - HTTP server port
- `GRPC_PORT` - gRPC server port
- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- `CORS_ALLOWED_ORIGINS` - Allowed CORS origins
- `MAX_CONN_POOL` - Database connection pool size
- `SHUTDOWN_TIMEOUT` - Graceful shutdown timeout

## API Endpoints

### HTTP API

#### Health Check
- GET `/health` - Service health status

#### Posts
- GET `/api/posts` - List posts (with pagination)
- GET `/api/posts/{id}` - Get post by ID
- POST `/api/posts` - Create new post
- PUT `/api/posts/{id}` - Update post
- DELETE `/api/posts/{id}` - Delete post

#### Replies
- GET `/api/posts/{id}/replies` - Get post replies
- POST `/api/posts/{id}/replies` - Create new reply
- DELETE `/api/posts/{id}/replies/{replyId}` - Delete reply

### gRPC API
The service also provides gRPC endpoints defined in `proto/forum.proto`.

## Running the Service

1. Set up environment variables in `.env` file
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the service:
   ```bash
   go run cmd/main.go
   ```

## API Request Examples

### Create Post
```http
POST /api/posts
Content-Type: application/json

{
    "title": "Example Post",
    "content": "Post content",
    "author_id": 1
}
```

### Create Reply
```http
POST /api/posts/1/replies
Content-Type: application/json

{
    "content": "Reply content",
    "author_id": 1
}
```

## Error Handling
The service returns standard HTTP status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

## Monitoring
Health check endpoint `/health` provides basic service status monitoring.