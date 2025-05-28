package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"backend/chat-service/internal/config"
	"backend/chat-service/internal/delivery/websocket"
	"backend/chat-service/internal/repository"
	"backend/chat-service/internal/usecase"
)

func main() {
	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Starting application with config",
		zap.String("db_host", cfg.Database.Host),
		zap.String("db_port", cfg.Database.Port),
		zap.String("db_name", cfg.Database.DBName),
		zap.String("db_user", cfg.Database.User),
		zap.String("http_port", cfg.Server.Port))

	// Подключение к базе данных
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Проверка подключения к базе данных
	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Successfully connected to database")

	// Проверка существования таблицы messages
	var exists bool
	err = pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'messages'
		)`).Scan(&exists)
	if err != nil {
		logger.Fatal("Failed to check if messages table exists", zap.Error(err))
	}
	if !exists {
		logger.Fatal("Messages table does not exist")
	}
	logger.Info("Messages table exists")

	// Инициализация слоев
	messageRepo := repository.NewMessageRepository(pool)
	chatUseCase := usecase.NewChatUseCase(messageRepo, pool)
	wsHandler := websocket.NewHandler(chatUseCase, logger)

	// Настройка маршрутизации
	router := mux.NewRouter()
	wsHandler.RegisterRoutes(router)

	// Запуск обработчика WebSocket соединений
	go chatUseCase.Run()

	// Настройка HTTP сервера
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-done
	logger.Info("Server stopping")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped gracefully")
}
