package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	wsHandler "backend/chat-service/internal/delivery/websocket"
	"backend/chat-service/internal/entity"
	"backend/chat-service/internal/repository"
	"backend/chat-service/internal/usecase"
)

var (
	testDBHost     = os.Getenv("TEST_DB_HOST")
	testDBPort     = os.Getenv("TEST_DB_PORT")
	testDBName     = os.Getenv("TEST_DB_NAME")
	testDBUser     = os.Getenv("TEST_DB_USER")
	testDBPassword = os.Getenv("TEST_DB_PASSWORD")
)

func setupTestDB(t *testing.T) (repository.MessageRepository, func()) {
	if testDBHost == "" {
		testDBHost = "localhost"
	}
	if testDBPort == "" {
		testDBPort = "5432"
	}
	if testDBName == "" {
		testDBName = "chat_test"
	}
	if testDBUser == "" {
		testDBUser = "postgres"
	}
	if testDBPassword == "" {
		testDBPassword = "postgres"
	}

	// Create database connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", testDBUser, testDBPassword, testDBHost, testDBPort, testDBName)

	// Create database connection
	pool, err := pgxpool.Connect(context.Background(), dbURL)
	require.NoError(t, err)

	// Create messages table
	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			user_id BIGINT NOT NULL,
			username VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS messages_created_at_idx ON messages(created_at DESC);
		CREATE INDEX IF NOT EXISTS messages_user_id_idx ON messages(user_id);
	`)
	require.NoError(t, err)

	repo := repository.NewMessageRepository(pool)

	// Return cleanup function
	cleanup := func() {
		// Clean up test data
		_, err := pool.Exec(context.Background(), "DELETE FROM messages")
		require.NoError(t, err)
		pool.Close()
	}

	return repo, cleanup
}

func TestWebSocketIntegration(t *testing.T) {
	// Skip if running in CI environment without DB
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Setup mock auth service
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/me" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       1,
				"username": "test_user",
			})
		}
	}))
	defer authServer.Close()

	// Override auth service URL in handler
	os.Setenv("AUTH_SERVICE_URL", authServer.URL)

	// Setup
	logger, _ := zap.NewDevelopment()
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	pool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", testDBUser, testDBPassword, testDBHost, testDBPort, testDBName))
	require.NoError(t, err)
	defer pool.Close()

	// Override WebSocket upgrader settings for tests
	wsHandler.SetUpgraderSettings(websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in tests
		},
	})

	useCase := usecase.NewChatUseCase(repo, pool)
	handler := wsHandler.NewHandler(useCase, logger)

	// Create test server
	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Test WebSocket Connection with Valid Token", func(t *testing.T) {
		// Create a valid token (in real scenario this would come from auth service)
		token := "valid_test_token"

		// Connect to WebSocket
		url := "ws" + server.URL[4:] + "/api/chat/ws?token=" + token
		c, resp, err := websocket.DefaultDialer.Dial(url, nil)
		require.NoError(t, err)
		defer c.Close()

		assert.Equal(t, 101, resp.StatusCode) // 101 is StatusSwitchProtocol

		// Wait for auth success message
		_, message, err := c.ReadMessage()
		require.NoError(t, err)

		var authResp usecase.ChatMessage
		err = json.Unmarshal(message, &authResp)
		require.NoError(t, err)
		assert.Equal(t, "auth_success", authResp.Type)
	})

	t.Run("Test Message History", func(t *testing.T) {
		// Insert test message
		ctx := context.Background()
		testMessage := entity.Message{
			ID:        "1",
			UserID:    1,
			Username:  "test_user",
			Content:   "test message",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, &testMessage)
		require.NoError(t, err)

		// Test history endpoint
		req := httptest.NewRequest("GET", "/api/chat/messages", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var messages []*entity.Message
		err = json.NewDecoder(w.Body).Decode(&messages)
		require.NoError(t, err)
		assert.NotEmpty(t, messages)
		assert.Equal(t, testMessage.Content, messages[0].Content)
	})
}
