package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiditAPIKey        string
	DiditWebhookSecret string
	DiditWebhookURL    string
	DiditWorkflowID    string
	DiditWorkflowURL   string // URL fixa do workflow
	ServerPort         string
	DatabaseURL        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		DiditAPIKey:        getEnv("DIDIT_API_KEY", ""),
		DiditWebhookSecret: getEnv("DIDIT_WEBHOOK_SECRET_KEY", ""),
		DiditWebhookURL:    getEnv("DIDIT_WEBHOOK_URL", ""),
		DiditWorkflowID:    getEnv("DIDIT_WORKFLOW_ID", ""),
		DiditWorkflowURL:   getEnv("DIDIT_WORKFLOW_SESSION", "https://verify.didit.me/verify/ynzf6TUsQ6BPDyVMWe8btQ"),
		ServerPort:         getEnv("SERVER_ADDR", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgresql://user:pass@localhost:5432/kyc_app?sslmode=disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
