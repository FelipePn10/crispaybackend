package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type VerificationRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type VerificationResponse struct {
	VerificationURL string `json:"verification_url"`
	UserID          string `json:"user_id"`
}

type WebhookEvent struct {
	EventType string          `json:"event_type"`
	Data      WebhookData     `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
	RawData   json.RawMessage `json:"raw_data,omitempty"` // To store the original payload
}

type WebhookData struct {
	SessionID string                 `json:"session_id"`
	Status    string                 `json:"status"`
	UserID    string                 `json:"user_id,omitempty"`
	UserData  map[string]interface{} `json:"user_data,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type VerificationSession struct {
	ID             uuid.UUID  `json:"id"`
	UserID         string     `json:"user_id"`
	SessionID      string     `json:"session_id"`
	DiditSessionID string     `json:"didit_session_id,omitempty"`
	Status         string     `json:"status"`
	UserEmail      string     `json:"user_email"`
	UserFirstName  string     `json:"user_first_name,omitempty"`
	UserLastName   string     `json:"user_last_name,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// WebhookEventDB represents the webhook event structure in the database.
type WebhookEventDB struct {
	ID        uuid.UUID       `json:"id"`
	EventType string          `json:"event_type"`
	SessionID string          `json:"session_id"`
	Payload   json.RawMessage `json:"payload"`
	Processed bool            `json:"processed"`
	CreatedAt time.Time       `json:"created_at"`
}
