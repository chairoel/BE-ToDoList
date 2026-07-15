package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBHost:    getEnv("POSTGRES_HOST", "localhost"),
		DBPort:    getEnv("POSTGRES_PORT", "5432"),
		DBUser:    getEnv("POSTGRES_USER", "postgres"),
		DBPass:    getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:    getEnv("POSTGRES_DB", "todolist"),
		DBSSLMode: getEnv("POSTGRES_SSLMODE", "disable"),
	}

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
