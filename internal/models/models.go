package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type VerificationRequest struct {
	UserId    string `json:"user_id" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Verificationresponse struct {
	VerificationURL string `json:"verification_url"`
	SessionID       string `json:"session_id"`
	Status          string `json:"status"`
}

type WebhookEvent struct {
	EventType string      `json:"event_type"`
	Data      WebhookData `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type WebhookData struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	UserID    string `json:"user_id"`
}

type VerificationSession struct {
	ID            uuid.UUID  `json:"id"`
	UserID        string     `json:"user_id"`
	SessionID     string     `json:"session_id"`
	DiditSession  string     `json:"didit_session,omitempty"`
	VeficatonURL  string     `json:"verification_url,omitempty"`
	Status        string     `json:"status"`
	UserEmail     string     `json:"user_email"`
	UserFisrtName string     `json:"user_first_name,omitempty"`
	UserLastName  string     `json:"user_last_name,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

// JSONB type for handling PostgreSQL JSONB
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return json.Unmarshal([]byte(v.(string)), j)
	}
	return json.Unmarshal(data, j)
}
