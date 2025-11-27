package config

import (
	"log"
	"syscall"

	"github.com/joho/godotenv"
)

type Config struct {
	DiditAPIKey        string
	DiditWebhookSecret string
	DiditWebhookURL    string
	DiditWorkflowID    string
	DiditBaseURL       string
	ServerPort         string
	DatabaseURL        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	return &Config{
		DiditAPIKey:        getEnv("DIDIT_API_KEY", ""),
		DiditWebhookSecret: getEnv("DIDIT_WEBHOOK_SECRET_KEY", ""),
		DiditWebhookURL:    getEnv("DIDIT_WEBHOOK_URL", ""),
		DiditWorkflowID:    getEnv("DIDIT_WORKFLOW_ID", ""),
		DiditBaseURL:       getEnv("DIDIT_BASE_URL", "https://api.didit.me"),
		ServerPort:         getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgresql://user:pass@localhost:5432/kyc?sslmode=disable"),
	}
}
func getEnv(key, defaultValue string) string {
	if value, exists := syscall.Getenv(key); exists {
		return value
	}
	return defaultValue
}
