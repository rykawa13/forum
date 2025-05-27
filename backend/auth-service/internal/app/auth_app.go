package app

import (
	"auth-service/internal/config"
	"auth-service/internal/delivery/router"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"auth-service/pkg/database"
	"auth-service/pkg/jwt"
	"fmt"
	"log"
)

type AuthApp struct {
	Config       *config.Config
	Repository   *repository.Repository
	UseCase      *usecase.AuthUseCase
	tokenManager jwt.TokenManager
}

func NewAuthApp(cfg *config.Config) *AuthApp {
	// Инициализация подключения к БД
	db, err := database.NewPostgresConnectionFromURL(cfg.DB.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Применение миграций
	if err := database.RunMigrations(cfg.DB.URL, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация репозиториев
	repos := repository.NewRepository(db)

	tokenManager, err := jwt.NewManager(cfg.JWT.Secret)
	if err != nil {
		log.Fatalf("Failed to create token manager: %v", err)
	}

	// Инициализация use case с конфигурацией
	authUC := usecase.NewAuthUseCase(
		repos.UserRepository,
		repos.SessionRepository,
		tokenManager,
		&usecase.Config{
			AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
			RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		},
	)

	return &AuthApp{
		Config:       cfg,
		Repository:   repos,
		UseCase:      authUC,
		tokenManager: tokenManager,
	}
}

func (a *AuthApp) Run() error {
	// Создаем роутер
	router := router.NewRouter(a.UseCase, a.tokenManager)

	// Запускаем сервер
	return router.Run(fmt.Sprintf(":%s", a.Config.Server.HTTPPort))
}
