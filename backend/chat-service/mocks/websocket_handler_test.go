package mocks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"backend/chat-service/internal/delivery/websocket"
	"backend/chat-service/internal/entity"
	"backend/chat-service/internal/repository"
	"backend/chat-service/internal/usecase"
)

// MockMessageRepository is a mock implementation of MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *entity.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetHistory(ctx context.Context, limit int32, beforeID int64) ([]*entity.Message, error) {
	args := m.Called(ctx, limit, beforeID)
	return args.Get(0).([]*entity.Message), args.Error(1)
}

func (m *MockMessageRepository) DeleteOldMessages(ctx context.Context, before time.Time) (int32, error) {
	args := m.Called(ctx, before)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockMessageRepository) QueryRow(ctx context.Context, query string, args ...interface{}) repository.Row {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(repository.Row)
}

func TestHandler_handleWebSocket(t *testing.T) {
	// Initialize logger
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		mockSetup      func(*MockMessageRepository)
	}{
		{
			name:           "No token provided",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			mockSetup:      func(m *MockMessageRepository) {},
		},
		{
			name:           "Invalid token",
			token:          "invalid_token",
			expectedStatus: http.StatusUnauthorized,
			mockSetup:      func(m *MockMessageRepository) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockMessageRepository)
			tt.mockSetup(mockRepo)

			chatUseCase := usecase.NewChatUseCase(mockRepo, nil)
			handler := websocket.NewHandler(chatUseCase, logger)

			// Create router and register routes
			router := mux.NewRouter()
			handler.RegisterRoutes(router)

			// Create test server
			server := httptest.NewServer(router)
			defer server.Close()

			// Create WebSocket URL
			url := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/chat/ws"
			if tt.token != "" {
				url += "?token=" + tt.token
			}

			// Try to connect
			_, resp, err := ws.DefaultDialer.Dial(url, nil)
			if resp.StatusCode != 101 { // 101 is StatusSwitchProtocol
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandler_handleGetHistory(t *testing.T) {
	// Initialize logger
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		limit          int32
		beforeID       int64
		expectedStatus int
		mockSetup      func(*MockMessageRepository)
		expectedBody   []*entity.Message
	}{
		{
			name:           "Successful history retrieval",
			limit:          50,
			beforeID:       0,
			expectedStatus: http.StatusOK,
			mockSetup: func(m *MockMessageRepository) {
				messages := []*entity.Message{
					{
						ID:        "1",
						UserID:    1,
						Username:  "test_user",
						Content:   "test message",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				m.On("GetHistory", mock.Anything, int32(50), int64(0)).Return(messages, nil)
			},
			expectedBody: []*entity.Message{
				{
					ID:       "1",
					UserID:   1,
					Username: "test_user",
					Content:  "test message",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockMessageRepository)
			tt.mockSetup(mockRepo)

			chatUseCase := usecase.NewChatUseCase(mockRepo, nil)
			handler := websocket.NewHandler(chatUseCase, logger)

			// Create request
			req := httptest.NewRequest("GET", "/api/chat/messages", nil)
			w := httptest.NewRecorder()

			// Create router and register routes
			router := mux.NewRouter()
			handler.RegisterRoutes(router)

			// Serve request
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []*entity.Message
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody[0].Content, response[0].Content)
				assert.Equal(t, tt.expectedBody[0].UserID, response[0].UserID)
				assert.Equal(t, tt.expectedBody[0].Username, response[0].Username)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
