package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"backend/chat-service/internal/entity"
	"backend/chat-service/internal/repository"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ChatMessage представляет сообщение для WebSocket
type ChatMessage struct {
	Type      string    `json:"type"`
	Content   string    `json:"content,omitempty"`
	UserID    int64     `json:"user_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Error     string    `json:"error,omitempty"`
	Token     string    `json:"token,omitempty"`
	ID        string    `json:"id,omitempty"`
	TempID    string    `json:"tempId,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Client представляет подключенного клиента
type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   int64
	Username string
	IsAuth   bool
	Role     string // "anonymous" или "authenticated"
	ctx      context.Context
	cancel   context.CancelFunc
}

// ChatUseCase представляет use case для чата
type ChatUseCase struct {
	repo       repository.MessageRepository
	db         *pgxpool.Pool
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// NewChatUseCase создает новый use case для чата
func NewChatUseCase(repo repository.MessageRepository, db *pgxpool.Pool) *ChatUseCase {
	return &ChatUseCase{
		repo:       repo,
		db:         db,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// NewClient создает нового клиента
func NewClient(conn *websocket.Conn, userID int64, username string, isAuth bool) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	role := "anonymous"
	if isAuth {
		role = "authenticated"
	}
	return &Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
		IsAuth:   isAuth,
		Role:     role,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Close закрывает клиента
func (c *Client) Close() {
	c.cancel()
	c.Conn.Close()
}

// Run запускает обработку WebSocket соединений
func (uc *ChatUseCase) Run() {
	for {
		select {
		case client := <-uc.Register:
			uc.mutex.Lock()
			uc.clients[client] = true
			uc.mutex.Unlock()

		case client := <-uc.unregister:
			uc.mutex.Lock()
			if _, ok := uc.clients[client]; ok {
				delete(uc.clients, client)
				close(client.Send)
			}
			uc.mutex.Unlock()

		case message := <-uc.broadcast:
			uc.mutex.RLock()
			for client := range uc.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(uc.clients, client)
				}
			}
			uc.mutex.RUnlock()
		}
	}
}

// HandleMessage обрабатывает входящее сообщение
func (c *Client) HandleMessage(message []byte, uc *ChatUseCase) error {
	var msg ChatMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}

	// Проверяем права на отправку сообщений
	if msg.Type == "message" && !c.IsAuth {
		errorMsg := ChatMessage{
			Type:  "error",
			Error: "Только авторизованные пользователи могут отправлять сообщения",
		}
		errorJSON, _ := json.Marshal(errorMsg)
		c.Send <- errorJSON
		return errors.New("unauthorized to send messages")
	}

	switch msg.Type {
	case "message":
		// Создаем новое сообщение
		newMsg := entity.Message{
			ID:        uuid.New().String(),
			Content:   msg.Content,
			UserID:    c.UserID,
			Username:  c.Username,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Сохраняем сообщение в БД
		if err := uc.repo.Create(c.ctx, &newMsg); err != nil {
			return fmt.Errorf("failed to save message: %v", err)
		}

		// Отправляем сообщение всем клиентам
		response := ChatMessage{
			Type:      "message",
			ID:        newMsg.ID,
			Content:   newMsg.Content,
			UserID:    newMsg.UserID,
			Username:  newMsg.Username,
			CreatedAt: newMsg.CreatedAt,
			TempID:    msg.TempID,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return fmt.Errorf("failed to marshal response: %v", err)
		}

		uc.broadcast <- responseJSON
	}

	return nil
}

// GetHistory возвращает историю сообщений
func (uc *ChatUseCase) GetHistory(ctx context.Context, limit int32, beforeID int64) ([]*entity.Message, error) {
	return uc.repo.GetHistory(ctx, limit, beforeID)
}

// DeleteOldMessages удаляет старые сообщения
func (uc *ChatUseCase) DeleteOldMessages(ctx context.Context, before time.Time) (int32, error) {
	return uc.repo.DeleteOldMessages(ctx, before)
}

// WritePump отправляет сообщения клиенту
func (c *Client) WritePump() {
	defer func() {
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// ReadPump читает сообщения от клиента
func (c *Client) ReadPump(uc *ChatUseCase) {
	defer func() {
		uc.unregister <- c
		c.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}

			if err := c.HandleMessage(message, uc); err != nil {
				log.Printf("error handling message: %v", err)
				errorMsg := ChatMessage{
					Type:  "error",
					Error: "Failed to process message",
				}
				errorData, _ := json.Marshal(errorMsg)
				c.Send <- errorData
			}
		}
	}
}
