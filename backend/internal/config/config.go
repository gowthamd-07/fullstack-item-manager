package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBUrl    string
	LogLevel string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists, but don't fail if it doesn't (k8s/docker might inject env vars directly)
	_ = godotenv.Load()

	cfg := &Config{
		Port:     getEnv("PORT", "8000"),
		DBUrl:    getEnv("DB_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
