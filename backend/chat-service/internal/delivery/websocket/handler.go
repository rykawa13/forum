package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"backend/chat-service/internal/usecase"
)

// @title Chat Service WebSocket API
// @version 1.0
// @description WebSocket API для чат-сервиса

// @host localhost:8080
// @BasePath /ws

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000" || // React dev server
			origin == "http://localhost:8080" || // Production frontend
			origin == "http://localhost:8081" || // Auth service
			origin == "http://localhost:8082" || // Forum service
			origin == "http://localhost:8083" // Chat service
	},
}

// SetUpgraderSettings allows overriding upgrader settings for tests
func SetUpgraderSettings(settings websocket.Upgrader) {
	upgrader = settings
}

// Handler обработчик WebSocket
type Handler struct {
	useCase *usecase.ChatUseCase
	logger  *zap.Logger
}

// NewHandler создает новый WebSocket обработчик
func NewHandler(useCase *usecase.ChatUseCase, logger *zap.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		logger:  logger,
	}
}

// RegisterRoutes регистрирует маршруты
func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Добавляем middleware для CORS и логирования
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.logger.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr))

			// Получаем origin из заголовка
			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = "*"
			}

			// Устанавливаем CORS заголовки
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// API endpoints
	api := r.PathPrefix("/api/chat").Subrouter()
	api.HandleFunc("/messages", h.handleGetHistory).Methods("GET", "OPTIONS")
	api.HandleFunc("/ws", h.handleWebSocket).Methods("GET", "OPTIONS") // WebSocket endpoint

	// Добавляем обработчик для проверки здоровья сервиса
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}

// @Summary Подключение к WebSocket чату
// @Description Устанавливает WebSocket соединение для обмена сообщениями
// @Tags websocket
// @Accept  json
// @Produce  json
// @Param   token     header    string     true        "Auth token"
// @Success 101 {object} usecase.ChatMessage
// @Router /chat [get]
func (h *Handler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("WebSocket connection attempt",
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("origin", r.Header.Get("Origin")))

	// Проверяем наличие токена
	token := r.URL.Query().Get("token")
	var authResp *struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}
	var err error

	// Если токен предоставлен, пытаемся аутентифицировать пользователя
	if token != "" {
		authResp, err = h.validateToken(r.Context(), token)
		if err != nil {
			h.logger.Warn("Failed to validate token",
				zap.Error(err),
				zap.String("remote_addr", r.RemoteAddr))
			// Продолжаем как анонимный пользователь
			authResp = nil
		}
	}

	// Устанавливаем WebSocket соединение
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection",
			zap.Error(err),
			zap.String("remote_addr", r.RemoteAddr))
		return
	}

	var client *usecase.Client
	if authResp != nil && authResp.ID != 0 && authResp.Username != "" {
		// Создаем аутентифицированного клиента
		client = usecase.NewClient(conn, authResp.ID, authResp.Username, true)
		h.logger.Info("Authenticated WebSocket connection established",
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int64("user_id", authResp.ID),
			zap.String("username", authResp.Username))

		// Отправляем подтверждение успешной аутентификации
		authSuccess := usecase.ChatMessage{
			Type:     "auth_success",
			UserID:   authResp.ID,
			Username: authResp.Username,
		}

		authSuccessJSON, err := json.Marshal(authSuccess)
		if err != nil {
			h.logger.Error("Failed to marshal auth success message",
				zap.Error(err),
				zap.Any("auth_success", authSuccess))
			conn.Close()
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, authSuccessJSON); err != nil {
			h.logger.Error("Failed to send auth success message",
				zap.Error(err),
				zap.Any("auth_success", authSuccess))
			conn.Close()
			return
		}
	} else {
		// Создаем анонимного клиента
		client = usecase.NewClient(conn, 0, "anonymous", false)
		h.logger.Info("Anonymous WebSocket connection established",
			zap.String("remote_addr", r.RemoteAddr))

		// Отправляем информацию о статусе анонимного пользователя
		anonInfo := usecase.ChatMessage{
			Type:  "connection_info",
			Error: "Вы подключены как анонимный пользователь. Для отправки сообщений необходима авторизация.",
		}

		anonInfoJSON, err := json.Marshal(anonInfo)
		if err != nil {
			h.logger.Error("Failed to marshal anonymous info message",
				zap.Error(err))
			conn.Close()
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, anonInfoJSON); err != nil {
			h.logger.Error("Failed to send anonymous info message",
				zap.Error(err))
			conn.Close()
			return
		}
	}

	// Регистрируем клиента
	h.useCase.Register <- client

	// Запускаем горутины для чтения и записи сообщений
	go client.WritePump()
	go client.ReadPump(h.useCase)
}

// validateToken validates the token with the auth service and returns user info
func (h *Handler) validateToken(ctx context.Context, token string) (*struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}, error) {
	// Get auth service URL from environment
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8081"
	}

	// Make request to auth service
	req, err := http.NewRequestWithContext(ctx, "GET", authServiceURL+"/api/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var result struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// @Summary Получение истории сообщений
// @Description Возвращает историю сообщений чата
// @Tags chat
// @Accept  json
// @Produce  json
// @Param   limit     query    int     true        "Limit"
// @Param   before_id query    int     false       "Before message ID"
// @Success 200 {array}  usecase.ChatMessage
// @Router /history [get]
func (h *Handler) handleGetHistory(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Getting chat history",
		zap.String("remote_addr", r.RemoteAddr))

	// Добавляем CORS заголовки
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	limit := int32(50) // Default limit
	beforeID := int64(0)

	messages, err := h.useCase.GetHistory(ctx, limit, beforeID)
	if err != nil {
		h.logger.Error("Failed to get chat history",
			zap.Error(err),
			zap.String("remote_addr", r.RemoteAddr))
		http.Error(w, "Failed to get chat history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		h.logger.Error("Failed to encode messages",
			zap.Error(err),
			zap.String("remote_addr", r.RemoteAddr))
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
}
