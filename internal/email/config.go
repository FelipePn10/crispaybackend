package email

import (
	"os"
)

type EmailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderPass  string
	SenderName  string
}

func LoadConfigFromEnv() EmailConfig {
	return EmailConfig{
		SMTPHost:    getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:    getEnv("SMTP_PORT", "587"),
		SenderEmail: os.Getenv("EMAIL_SENDER"),
		SenderPass:  os.Getenv("EMAIL_PASSWORD"),
		SenderName:  getEnv("EMAIL_SENDER_NAME", "Sua Empresa"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
