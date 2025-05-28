package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config представляет конфигурацию приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// ServerConfig представляет конфигурацию сервера
type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

// DatabaseConfig представляет конфигурацию базы данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AuthConfig представляет конфигурацию аутентификации
type AuthConfig struct {
	AuthServiceURL string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Игнорируем ошибку, если .env файл не найден
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	shutdownTimeout, _ := strconv.Atoi(getEnv("SHUTDOWN_TIMEOUT_SECONDS", "5"))

	return &Config{
		Server: ServerConfig{
			Port:            getEnv("HTTP_PORT", "8080"),
			ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "chat"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		},
	}, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetDSN возвращает строку подключения к базе данных
func (c *DatabaseConfig) GetDSN() string {
	return "postgres://" + c.User + ":" + c.Password +
		"@" + c.Host + ":" + c.Port +
		"/" + c.DBName + "?sslmode=" + c.SSLMode
}
