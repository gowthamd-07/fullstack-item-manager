package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBUrl          string
	LogLevel       string
	AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	origins := getEnv("ALLOWED_ORIGINS", "http://localhost:8080")
	var parsedOrigins []string
	for _, o := range strings.Split(origins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			parsedOrigins = append(parsedOrigins, trimmed)
		}
	}

	cfg := &Config{
		Port:           getEnv("PORT", "8000"),
		DBUrl:          getEnv("DB_URL", ""),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		AllowedOrigins: parsedOrigins,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
