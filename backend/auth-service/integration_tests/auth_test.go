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
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"backend/auth-service/internal/delivery/http/handler"
	"backend/auth-service/internal/repository"
	"backend/auth-service/internal/usecase"
)

var (
	testDBHost     = os.Getenv("TEST_DB_HOST")
	testDBPort     = os.Getenv("TEST_DB_PORT")
	testDBName     = os.Getenv("TEST_DB_NAME")
	testDBUser     = os.Getenv("TEST_DB_USER")
	testDBPassword = os.Getenv("TEST_DB_PASSWORD")
)

func setupTestDB(t *testing.T) (repository.UserRepository, func()) {
	if testDBHost == "" {
		testDBHost = "localhost"
	}
	if testDBPort == "" {
		testDBPort = "5432"
	}
	if testDBName == "" {
		testDBName = "auth_test"
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
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	require.NoError(t, err)

	repo := repository.NewUserRepository(pool)

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "DROP TABLE IF EXISTS users")
		require.NoError(t, err)
		pool.Close()
	}

	return repo, cleanup
}

func TestAuthIntegration(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	logger, _ := zap.NewDevelopment()
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	useCase := usecase.NewAuthUseCase(repo, []byte("test_secret"), 24*time.Hour)
	handler := handler.NewAuthHandler(useCase, logger)

	t.Run("Full Authentication Flow", func(t *testing.T) {
		// 1. Register a new user
		registerInput := map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(registerInput)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// 2. Try to login with wrong password
		loginWrongInput := map[string]interface{}{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		body, _ = json.Marshal(loginWrongInput)
		req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handler.Login(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// 3. Login with correct credentials
		loginInput := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ = json.Marshal(loginInput)
		req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handler.Login(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var loginResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &loginResponse)
		token, ok := loginResponse["token"].(string)
		require.True(t, ok)
		require.NotEmpty(t, token)

		// 4. Validate token
		req = httptest.NewRequest("GET", "/api/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		handler.Me(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var userResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &userResponse)
		assert.Equal(t, "testuser", userResponse["username"])
		assert.Equal(t, "test@example.com", userResponse["email"])
	})

	t.Run("Registration Validation", func(t *testing.T) {
		tests := []struct {
			name         string
			input        map[string]interface{}
			expectedCode int
		}{
			{
				name: "Empty username",
				input: map[string]interface{}{
					"username": "",
					"email":    "test2@example.com",
					"password": "password123",
				},
				expectedCode: http.StatusBadRequest,
			},
			{
				name: "Invalid email",
				input: map[string]interface{}{
					"username": "testuser2",
					"email":    "invalid-email",
					"password": "password123",
				},
				expectedCode: http.StatusBadRequest,
			},
			{
				name: "Short password",
				input: map[string]interface{}{
					"username": "testuser2",
					"email":    "test2@example.com",
					"password": "short",
				},
				expectedCode: http.StatusBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.input)
				req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				handler.Register(w, req)
				assert.Equal(t, tt.expectedCode, w.Code)
			})
		}
	})
}
