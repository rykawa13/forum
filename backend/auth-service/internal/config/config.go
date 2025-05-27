package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	JWT      JWTConfig
	LogLevel string
}

type ServerConfig struct {
	HTTPPort string
	GRPCPort string
}

type DBConfig struct {
	URL string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			HTTPPort: getEnv("HTTP_PORT", "8081"),
			GRPCPort: getEnv("GRPC_PORT", "50051"),
		},
		DB: DBConfig{
			URL: getEnv("DB_URL", "postgres://forum_user:forum_password@localhost:5432/forum_db?sslmode=disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenTTL:  24 * time.Hour,
			RefreshTokenTTL: 30 * 24 * time.Hour,
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
