package mocks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"backend/auth-service/internal/entity"
	"backend/auth-service/internal/usecase"
)

// MockAuthUseCase мок для auth use case
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthUseCase) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUseCase) ValidateToken(token string) (*entity.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name         string
		input        map[string]interface{}
		setupMock    func(*MockAuthUseCase)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "Successful registration",
			input: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(m *MockAuthUseCase) {
				m.On("Register", mock.AnythingOfType("*entity.User")).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "User registered successfully",
			},
		},
		{
			name: "Invalid input",
			input: map[string]interface{}{
				"username": "",
				"email":    "invalid-email",
				"password": "short",
			},
			setupMock:    func(m *MockAuthUseCase) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid input data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockAuthUseCase)
			tt.setupMock(mockUseCase)

			handler := NewAuthHandler(mockUseCase, logger)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name          string
		input         map[string]interface{}
		setupMock     func(*MockAuthUseCase)
		expectedCode  int
		expectedToken bool
	}{
		{
			name: "Successful login",
			input: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(m *MockAuthUseCase) {
				m.On("Login", "test@example.com", "password123").Return("valid_token", nil)
			},
			expectedCode:  http.StatusOK,
			expectedToken: true,
		},
		{
			name: "Invalid credentials",
			input: map[string]interface{}{
				"email":    "wrong@example.com",
				"password": "wrongpass",
			},
			setupMock: func(m *MockAuthUseCase) {
				m.On("Login", "wrong@example.com", "wrongpass").Return("", usecase.ErrInvalidCredentials)
			},
			expectedCode:  http.StatusUnauthorized,
			expectedToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockAuthUseCase)
			tt.setupMock(mockUseCase)

			handler := NewAuthHandler(mockUseCase, logger)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedToken {
				assert.Contains(t, response, "token")
				assert.NotEmpty(t, response["token"])
			} else {
				assert.Contains(t, response, "error")
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}
