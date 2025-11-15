package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

func Load() (*Config, error) {
	godotenv.Load()

	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "meobeo-talk-api"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "123"),
			DBName:   getEnv("DB_NAME", "meobeo_talk"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			ExpireTime: 24 * time.Hour,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
