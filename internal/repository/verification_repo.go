package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FelipePn10/crispaybackend/internal/database/sqlc"
	"github.com/FelipePn10/crispaybackend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type VerificationRepository struct {
	queries *sqlc.Queries
}

func NewVerificationRepository(queries *sqlc.Queries) *VerificationRepository {
	return &VerificationRepository{
		queries: queries,
	}
}

func (r *VerificationRepository) CreateSession(ctx context.Context, params *models.VerificationSession) (*models.VerificationSession, error) {
	dbParams := sqlc.CreateVerificationSessionParams{
		UserID:         params.UserID,
		SessionID:      params.SessionID,
		DiditSessionID: sql.NullString{String: params.DiditSessionID, Valid: params.DiditSessionID != ""},
		UserEmail:      params.UserEmail,
		UserFirstName:  sql.NullString{String: params.UserFirstName, Valid: params.UserFirstName != ""},
		UserLastName:   sql.NullString{String: params.UserLastName, Valid: params.UserLastName != ""},
		Status:         params.Status,
	}

	result, err := r.queries.CreateVerificationSession(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create verification session: %v", err)
	}

	return r.toDomainModel(result), nil
}

func (r *VerificationRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*models.VerificationSession, error) {
	result, err := r.queries.GetVerificationSessionByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get verification session: %v", err)
	}

	return r.toDomainModel(result), nil
}

func (r *VerificationRepository) GetSessionBySessionID(ctx context.Context, sessionID string) (*models.VerificationSession, error) {
	result, err := r.queries.GetVerificationSessionBySessionID(ctx, sessionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get verification session: %v", err)
	}
	return r.toDomainModel(result), nil
}

func (r *VerificationRepository) GetSessionByDiditSessionID(ctx context.Context, diditSessionID string) (*models.VerificationSession, error) {
	nullDiditSessionID := sql.NullString{String: diditSessionID, Valid: diditSessionID != ""}

	result, err := r.queries.GetVerificationSessionByDiditSessionID(ctx, nullDiditSessionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get verification session by didit session id: %v", err)
	}

	return r.toDomainModel(result), nil
}

func (r *VerificationRepository) UpdateStatus(ctx context.Context, sessionID string, status string) (*models.VerificationSession, error) {
	result, err := r.queries.UpdateVerificationSessionStatus(ctx, sqlc.UpdateVerificationSessionStatusParams{
		SessionID: sessionID,
		Status:    status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update verification session status: %v", err)
	}
	return r.toDomainModel(result), nil
}

func (r *VerificationRepository) UpdateDiditSessionID(ctx context.Context, sessionID string, diditSessionID string) (*models.VerificationSession, error) {
	result, err := r.queries.UpdateDiditSessionID(ctx, sqlc.UpdateDiditSessionIDParams{
		SessionID:      sessionID,
		DiditSessionID: sql.NullString{String: diditSessionID, Valid: diditSessionID != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update didit session id: %v", err)
	}
	return r.toDomainModel(result), nil
}

// UpdateDiditData is maintained for compatibility (it calls UpdateDiditSessionID internally)
func (r *VerificationRepository) UpdateDiditData(ctx context.Context, sessionID string, diditSessionID string, verificationURL string) (*models.VerificationSession, error) {
	return r.UpdateDiditSessionID(ctx, sessionID, diditSessionID)
}

func (r *VerificationRepository) ListVerificationSessionsByUserID(ctx context.Context, userID string) ([]*models.VerificationSession, error) {
	results, err := r.queries.ListVerificationSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list verification sessions: %v", err)
	}

	sessions := make([]*models.VerificationSession, len(results))
	for i, result := range results {
		sessions[i] = r.toDomainModel(result)
	}

	return sessions, nil
}

func (r *VerificationRepository) ListVerificationSessionsByStatus(ctx context.Context, status string) ([]*models.VerificationSession, error) {
	results, err := r.queries.ListVerificationSessionsByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list verification sessions by status: %v", err)
	}

	sessions := make([]*models.VerificationSession, len(results))
	for i, result := range results {
		sessions[i] = r.toDomainModel(result)
	}

	return sessions, nil
}

func (r *VerificationRepository) CreateWebhookEvent(ctx context.Context, eventType string, sessionID string, payload []byte) error {
	_, err := r.queries.CreateWebhookEvent(ctx, sqlc.CreateWebhookEventParams{
		EventType: eventType,
		SessionID: sessionID,
		Payload:   payload,
	})
	if err != nil {
		return fmt.Errorf("failed to create webhook event: %v", err)
	}
	return nil
}

func (r *VerificationRepository) GetWebhookEventsBySessionID(ctx context.Context, sessionID string) ([]*models.WebhookEvent, error) {
	results, err := r.queries.GetWebhookEventsBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook events: %v", err)
	}

	events := make([]*models.WebhookEvent, len(results))
	for i, result := range results {
		event, err := r.webhookEventToDomainModel(result)
		if err != nil {
			return nil, fmt.Errorf("failed to convert webhook event: %v", err)
		}
		events[i] = event
	}

	return events, nil
}

func (r *VerificationRepository) toDomainModel(dbSession sqlc.VerificationSession) *models.VerificationSession {
	session := &models.VerificationSession{
		ID:        dbSession.ID,
		UserID:    dbSession.UserID,
		SessionID: dbSession.SessionID,
		Status:    dbSession.Status,
		UserEmail: dbSession.UserEmail,
		CreatedAt: dbSession.CreatedAt,
		UpdatedAt: dbSession.UpdatedAt,
	}

	if dbSession.DiditSessionID.Valid {
		session.DiditSessionID = dbSession.DiditSessionID.String
	}
	if dbSession.UserFirstName.Valid {
		session.UserFirstName = dbSession.UserFirstName.String
	}
	if dbSession.UserLastName.Valid {
		session.UserLastName = dbSession.UserLastName.String
	}
	if dbSession.CompletedAt.Valid {
		session.CompletedAt = &dbSession.CompletedAt.Time
	}

	return session
}

func (r *VerificationRepository) webhookEventToDomainModel(dbWebhook sqlc.WebhookEvent) (*models.WebhookEvent, error) {
	var webhookEvent models.WebhookEvent

	// Tentar fazer parse do payload como WebhookEvent primeiro
	if err := json.Unmarshal(dbWebhook.Payload, &webhookEvent); err != nil {
		// Se não der, tentar fazer parse manualmente
		var rawData map[string]interface{}
		if err := json.Unmarshal(dbWebhook.Payload, &rawData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal webhook payload: %v", err)
		}

		// Construir o WebhookEvent manualmente
		webhookEvent = models.WebhookEvent{
			RawData: dbWebhook.Payload,
		}

		// Extrair event_type
		if eventType, ok := rawData["event_type"].(string); ok {
			webhookEvent.EventType = eventType
		} else if eventType, ok := rawData["type"].(string); ok {
			webhookEvent.EventType = eventType
		}

		// Extrair timestamp
		if timestampStr, ok := rawData["timestamp"].(string); ok {
			if timestamp, err := time.Parse(time.RFC3339, timestampStr); err == nil {
				webhookEvent.Timestamp = timestamp
			}
		} else if timestamp, ok := rawData["timestamp"].(float64); ok {
			// Se for unix timestamp
			webhookEvent.Timestamp = time.Unix(int64(timestamp), 0)
		} else {
			// Usar created_at do banco como fallback
			webhookEvent.Timestamp = dbWebhook.CreatedAt
		}

		// Extrair data
		webhookEvent.Data = models.WebhookData{}

		if data, ok := rawData["data"].(map[string]interface{}); ok {
			if sessionID, ok := data["session_id"].(string); ok {
				webhookEvent.Data.SessionID = sessionID
			}

			if status, ok := data["status"].(string); ok {
				webhookEvent.Data.Status = status
			}

			if userID, ok := data["user_id"].(string); ok {
				webhookEvent.Data.UserID = userID
			}

			if userData, ok := data["user_data"].(map[string]interface{}); ok {
				webhookEvent.Data.UserData = userData
			}

			if metadata, ok := data["metadata"].(map[string]interface{}); ok {
				webhookEvent.Data.Metadata = metadata
			}
		} else {
			// Se não houver estrutura "data", tentar extrair diretamente do root
			if sessionID, ok := rawData["session_id"].(string); ok {
				webhookEvent.Data.SessionID = sessionID
			}

			if status, ok := rawData["status"].(string); ok {
				webhookEvent.Data.Status = status
			}

			if userID, ok := rawData["user_id"].(string); ok {
				webhookEvent.Data.UserID = userID
			}
		}
	} else {
		// Se o parse direto funcionou, armazenar o raw data também
		webhookEvent.RawData = dbWebhook.Payload
	}

	return &webhookEvent, nil
}
