package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBSource      string        `json:"db_source"`
	ServerAddress string        `json:"server_address"`
	JWTSecret     string        `json:"jwt_secret"`
	TokenDuration time.Duration `json:"token_duration"`
}

func LoadConfig(path string) (config Config, err error) {
	envFile := filepath.Join(path, ".env")
	if _, statErr := os.Stat(envFile); statErr == nil {
		if loadErr := godotenv.Load(envFile); loadErr != nil {
			return config, fmt.Errorf("failed to load .env file: %w", loadErr)
		}
	}

	config.DBSource = getEnv("DB_SOURCE", "postgresql://root:secret@localhost:5432/bank_ledger?sslmode=disable")
	config.ServerAddress = getEnv("SERVER_ADDRESS", "0.0.0.0:8080")
	config.JWTSecret = getEnv("JWT_SECRET", "default_secret_key_change_me")

	durationStr := getEnv("JWT_DURATION", "15m")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return config, fmt.Errorf("invalid JWT_DURATION format (%s): %w", durationStr, err)
	}
	config.TokenDuration = duration

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
