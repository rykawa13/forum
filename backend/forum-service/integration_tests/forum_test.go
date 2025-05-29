package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"backend/forum-service/internal/delivery/http/handler"
	"backend/forum-service/internal/repository"
	"backend/forum-service/internal/usecase"
)

var (
	testDBHost     = os.Getenv("TEST_DB_HOST")
	testDBPort     = os.Getenv("TEST_DB_PORT")
	testDBName     = os.Getenv("TEST_DB_NAME")
	testDBUser     = os.Getenv("TEST_DB_USER")
	testDBPassword = os.Getenv("TEST_DB_PASSWORD")
)

func setupTestDB(t *testing.T) (repository.ForumRepository, func()) {
	if testDBHost == "" {
		testDBHost = "localhost"
	}
	if testDBPort == "" {
		testDBPort = "5432"
	}
	if testDBName == "" {
		testDBName = "forum_test"
	}
	if testDBUser == "" {
		testDBUser = "postgres"
	}
	if testDBPassword == "" {
		testDBPassword = "postgres"
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", testDBUser, testDBPassword, testDBHost, testDBPort, testDBName)

	pool, err := pgxpool.Connect(context.Background(), dbURL)
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS topics (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id BIGINT NOT NULL,
			username TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			user_id BIGINT NOT NULL,
			username TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS topics_created_at_idx ON topics(created_at DESC);
		CREATE INDEX IF NOT EXISTS posts_topic_id_idx ON posts(topic_id);
		CREATE INDEX IF NOT EXISTS posts_created_at_idx ON posts(created_at DESC);
	`)
	require.NoError(t, err)

	repo := repository.NewForumRepository(pool)

	cleanup := func() {
		_, err := pool.Exec(context.Background(), `
			DROP TABLE IF EXISTS posts;
			DROP TABLE IF EXISTS topics;
		`)
		require.NoError(t, err)
		pool.Close()
	}

	return repo, cleanup
}

func TestForumIntegration(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	logger, _ := zap.NewDevelopment()
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	useCase := usecase.NewForumUseCase(repo)
	handler := handler.NewForumHandler(useCase, logger)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	t.Run("Full Forum Flow", func(t *testing.T) {
		// 1. Create a topic
		topicInput := map[string]interface{}{
			"title":   "Test Topic",
			"content": "Test topic content",
		}
		body, _ := json.Marshal(topicInput)
		req := httptest.NewRequest("POST", "/api/forum/topics", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")
		req.Header.Set("X-Username", "testuser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var topicResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &topicResponse)
		topicID := topicResponse["id"].(string)
		require.NotEmpty(t, topicID)

		// 2. Get the created topic
		req = httptest.NewRequest("GET", "/api/forum/topics/"+topicID, nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var topic map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &topic)
		assert.Equal(t, "Test Topic", topic["title"])
		assert.Equal(t, "Test topic content", topic["content"])

		// 3. Create a post in the topic
		postInput := map[string]interface{}{
			"topic_id": topicID,
			"content":  "Test post content",
		}
		body, _ = json.Marshal(postInput)
		req = httptest.NewRequest("POST", "/api/forum/posts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")
		req.Header.Set("X-Username", "testuser")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// 4. Get topic with posts
		req = httptest.NewRequest("GET", "/api/forum/topics/"+topicID+"/posts", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var postsResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &postsResponse)
		posts := postsResponse["posts"].([]interface{})
		assert.Len(t, posts, 1)
		assert.Equal(t, "Test post content", posts[0].(map[string]interface{})["content"])
	})

	t.Run("List Topics", func(t *testing.T) {
		// Create multiple topics
		for i := 1; i <= 3; i++ {
			topicInput := map[string]interface{}{
				"title":   fmt.Sprintf("Topic %d", i),
				"content": fmt.Sprintf("Content %d", i),
			}
			body, _ := json.Marshal(topicInput)
			req := httptest.NewRequest("POST", "/api/forum/topics", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", "1")
			req.Header.Set("X-Username", "testuser")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusCreated, w.Code)
		}

		// List topics
		req := httptest.NewRequest("GET", "/api/forum/topics?page=1&limit=10", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		topics := response["topics"].([]interface{})
		assert.GreaterOrEqual(t, len(topics), 3)
	})

	t.Run("Topic Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/forum/topics/non-existent-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		// Try to create topic without required fields
		topicInput := map[string]interface{}{
			"title": "",
		}
		body, _ := json.Marshal(topicInput)
		req := httptest.NewRequest("POST", "/api/forum/topics", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1")
		req.Header.Set("X-Username", "testuser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
