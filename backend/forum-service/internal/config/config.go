package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Logger LoggerConfig
	HTTP   HTTPConfig
	GRPC   GRPCConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	URL string
}

type LoggerConfig struct {
	Level string
}

type HTTPConfig struct {
	Port string
}

type GRPCConfig struct {
	Port int
}

func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8082"),
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		DB: DBConfig{
			URL: getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/forum?sslmode=disable"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "5000"),
		},
		GRPC: GRPCConfig{
			Port: getEnvInt("GRPC_PORT", 50052),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
