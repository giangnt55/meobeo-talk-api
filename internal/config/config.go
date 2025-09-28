package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// App
	AppName  string
	AppEnv   string
	AppPort  string
	AppDebug bool

	// Database
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string
	DBSSLMode  string
	DBTimezone string

	// JWT
	JWTSecret      string
	JWTExpireHours int
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		// Not fatal in production where env vars are set differently
		if os.Getenv("APP_ENV") != "production" {
			return nil, err
		}
	}

	return &Config{
		AppName:  getEnv("APP_NAME", "MyAPI"),
		AppEnv:   getEnv("APP_ENV", "development"),
		AppPort:  getEnv("APP_PORT", "8080"),
		AppDebug: getEnvBool("APP_DEBUG", true),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBDatabase: getEnv("DB_DATABASE", "myapp_db"),
		DBUsername: getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		DBTimezone: getEnv("DB_TIMEZONE", "UTC"),

		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "true" {
		return true
	} else if value == "false" {
		return false
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	// Implementation for parsing int from env
	return defaultValue
}
