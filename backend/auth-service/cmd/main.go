package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	// Создаем конфигурацию
	cfg := config.New()

	// Создаем и запускаем приложение
	application := app.NewAuthApp(cfg)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
