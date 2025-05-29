package mocks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"backend/forum-service/internal/entity"
)

// MockForumUseCase мок для forum use case
type MockForumUseCase struct {
	mock.Mock
}

func (m *MockForumUseCase) CreateTopic(topic *entity.Topic) error {
	args := m.Called(topic)
	return args.Error(0)
}

func (m *MockForumUseCase) GetTopic(id string) (*entity.Topic, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Topic), args.Error(1)
}

func (m *MockForumUseCase) ListTopics(page, limit int) ([]*entity.Topic, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]*entity.Topic), args.Error(1)
}

func (m *MockForumUseCase) CreatePost(post *entity.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func TestForumHandler_CreateTopic(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name         string
		input        map[string]interface{}
		setupAuth    func(*http.Request)
		setupMock    func(*MockForumUseCase)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "Successful topic creation",
			input: map[string]interface{}{
				"title":   "Test Topic",
				"content": "Test content",
			},
			setupAuth: func(r *http.Request) {
				r.Header.Set("X-User-ID", "1")
				r.Header.Set("X-Username", "testuser")
			},
			setupMock: func(m *MockForumUseCase) {
				m.On("CreateTopic", mock.AnythingOfType("*entity.Topic")).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Topic created successfully",
			},
		},
		{
			name: "Unauthorized request",
			input: map[string]interface{}{
				"title":   "Test Topic",
				"content": "Test content",
			},
			setupAuth:    func(r *http.Request) {},
			setupMock:    func(m *MockForumUseCase) {},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Unauthorized",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockForumUseCase)
			tt.setupMock(mockUseCase)

			handler := NewForumHandler(mockUseCase, logger)
			router := mux.NewRouter()
			handler.RegisterRoutes(router)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/api/forum/topics", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			tt.setupAuth(req)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestForumHandler_GetTopic(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name          string
		topicID       string
		setupMock     func(*MockForumUseCase)
		expectedCode  int
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "Existing topic",
			topicID: "1",
			setupMock: func(m *MockForumUseCase) {
				m.On("GetTopic", "1").Return(&entity.Topic{
					ID:      "1",
					Title:   "Test Topic",
					Content: "Test content",
					UserID:  1,
				}, nil)
			},
			expectedCode: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, "Test Topic", response["title"])
			},
		},
		{
			name:    "Non-existing topic",
			topicID: "999",
			setupMock: func(m *MockForumUseCase) {
				m.On("GetTopic", "999").Return(nil, entity.ErrTopicNotFound)
			},
			expectedCode: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockForumUseCase)
			tt.setupMock(mockUseCase)

			handler := NewForumHandler(mockUseCase, logger)
			router := mux.NewRouter()
			handler.RegisterRoutes(router)

			req := httptest.NewRequest("GET", "/api/forum/topics/"+tt.topicID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			tt.checkResponse(t, w)
			mockUseCase.AssertExpectations(t)
		})
	}
}
