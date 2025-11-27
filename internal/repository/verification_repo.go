package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FelipePn10/crispaybackend/internal/database/sqlc"
	"github.com/FelipePn10/crispaybackend/internal/models"
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

func (r *VerificationRepository) CreateSession(ctx context.Context, params models.VerificationSession) (*models.VerificationSession, error) {
	dbParams := sqlc.CreateVerificationSessionParams{
		UserID:    params.UserID,
		SessionID: params.SessionID,
		DiditSessionID: sql.NullString{
			String: params.DiditSession,
			Valid:  params.DiditSession != "",
		},
		VerificationUrl: sql.NullString{
			String: params.VeficatonURL,
			Valid:  params.VeficatonURL != "",
		},
		UserEmail: params.UserEmail,
		UserFirstName: sql.NullString{
			String: params.UserFisrtName,
			Valid:  params.UserFisrtName != "",
		},
		UserLastName: sql.NullString{
			String: params.UserLastName,
			Valid:  params.UserLastName != "",
		},
		Status: params.Status,
	}
	result, err := r.queries.CreateVerificationSession(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create verification session: %v", err)
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

func (r *VerificationRepository) UpdateDiditData(ctx context.Context, sessionID string, diditSessionID string, verificationURL string) (*models.VerificationSession, error) {
	result, err := r.queries.UpdateDiditSessionData(ctx, sqlc.UpdateDiditSessionDataParams{
		SessionID:       sessionID,
		DiditSessionID:  sql.NullString{String: diditSessionID, Valid: diditSessionID != ""},
		VerificationUrl: sql.NullString{String: verificationURL, Valid: verificationURL != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update verification session didit data: %v", err)
	}
	return r.toDomainModel(result), nil
}

func (r *VerificationRepository) CreateWebhookEvent(ctx context.Context, eventType string, sessionId string, payload []byte) error {
	_, err := r.queries.CreateWebhookEvent(ctx, sqlc.CreateWebhookEventParams{
		EventType: eventType,
		SessionID: sessionId,
		Payload:   payload,
	})
	if err != nil {
		return fmt.Errorf("failed to create webhook event: %v", err)
	}
	return nil
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
		session.DiditSession = dbSession.DiditSessionID.String
	}
	if dbSession.VerificationUrl.Valid {
		session.VeficatonURL = dbSession.VerificationUrl.String
	}
	if dbSession.UserFirstName.Valid {
		session.UserFisrtName = dbSession.UserFirstName.String
	}
	if dbSession.UserLastName.Valid {
		session.UserLastName = dbSession.UserLastName.String
	}
	if dbSession.CompletedAt.Valid {
		session.CompletedAt = &dbSession.CompletedAt.Time
	}

	return session
}
