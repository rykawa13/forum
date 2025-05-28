package app

import (
	"context"
	"forum-service/internal/config"
	delivery "forum-service/internal/delivery/http"
	"forum-service/internal/logger"
	"forum-service/internal/repository"
	"forum-service/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func NewServer(postUC usecase.PostUseCase) *Server {
	handler := delivery.NewHandler(postUC)

	srv := &http.Server{
		Handler:      handler.Init(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	logger, _ := logger.NewLogger("debug")

	return &Server{
		httpServer: srv,
		logger:     logger,
	}
}

func (s *Server) Start(addr string) error {
	s.httpServer.Addr = addr
	s.logger.Info("Starting server", zap.String("address", addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// App represents the application
type App struct {
	config     *config.Config
	logger     *zap.Logger
	httpServer *http.Server
}

// New creates a new application instance
func New(configPath string) (*App, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.Logger.Level)
	if err != nil {
		return nil, err
	}

	// Connect to database
	pool, err := pgxpool.Connect(context.Background(), cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	repo := repository.NewPostRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	// Initialize use cases
	uc := usecase.NewPostUseCase(repo, userRepo)

	// Initialize HTTP server
	handler := delivery.NewHandler(uc)
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler.Init(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &App{
		config:     cfg,
		logger:     log,
		httpServer: srv,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	// Channel for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server
	go func() {
		a.logger.Info("Starting HTTP server", zap.String("port", a.config.Server.Port))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-quit
	a.logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	a.logger.Info("Server exited properly")
	return nil
}
