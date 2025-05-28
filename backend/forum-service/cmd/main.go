package main

import (
	"context"
	"fmt"
	"forum-service/internal/app"
	"forum-service/internal/config"
	"forum-service/internal/logger"
	"forum-service/internal/repository"
	"forum-service/internal/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found or error loading it: %v\n", err)
	}

	// Initialize logger
	zapLogger, err := logger.NewLogger("info")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer zapLogger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	logger.Info("Configuration loaded")

	// Database connection
	dbpool, err := pgxpool.Connect(context.Background(), cfg.DB.URL)
	if err != nil {
		logger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer dbpool.Close()
	logger.Info("Connected to database")

	// Initialize database schema
	if err := repository.InitDatabase(context.Background(), dbpool); err != nil {
		logger.Fatal("Failed to initialize database schema", zap.Error(err))
	}
	logger.Info("Database schema initialized")

	// Repository and use case initialization
	postRepo := repository.NewPostRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)
	postUseCase := usecase.NewPostUseCase(postRepo, userRepo)
	logger.Info("Initialized repository and use case")

	// Create and start HTTP server
	server := app.NewServer(postUseCase)
	go func() {
		if err := server.Start(fmt.Sprintf(":%s", cfg.HTTP.Port)); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
