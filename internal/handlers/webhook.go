package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/FelipePn10/crispaybackend/config"
	"github.com/FelipePn10/crispaybackend/internal/didit"
	"github.com/FelipePn10/crispaybackend/internal/email/service"
	"github.com/FelipePn10/crispaybackend/internal/models"
	"github.com/FelipePn10/crispaybackend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WebhookHandler struct {
	diditClient *didit.Client
	config      *config.Config
	repo        *repository.VerificationRepository
	email       *service.EmailService
}

func NewWebhookHandler(diditClient *didit.Client, cfg *config.Config, repo *repository.VerificationRepository, emailService *service.EmailService) *WebhookHandler {
	return &WebhookHandler{
		diditClient: diditClient,
		config:      cfg,
		repo:        repo,
		email:       emailService,
	}
}

// HandleVerificationWebhook processes Didit webhooks
func (h *WebhookHandler) HandleVerificationWebhook(c *gin.Context) {
	//signature := c.GetHeader("X-Didit-Signature")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Webhook received: %s", string(body))

	// Validate signature
	// if !h.diditClient.VerifyWebhookSignature(body, signature) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
	// 	return
	// }

	var webhookEvent models.WebhookEvent
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Save webhook to database
	if err := h.repo.CreateWebhookEvent(c.Request.Context(), webhookEvent.EventType, webhookEvent.Data.SessionID, body); err != nil {
		log.Printf("Failed to save webhook event: %v", err)
	}

	// Process events
	switch webhookEvent.EventType {
	case "verification.completed", "verification.approved":
		h.handleVerificationCompleted(c.Request.Context(), webhookEvent)
	case "verification.failed", "verification.rejected":
		h.handleVerificationFailed(c.Request.Context(), webhookEvent)
	case "verification.review":
		h.handleVerificationReview(c.Request.Context(), webhookEvent)
	default:
		log.Printf("Unhandled event type: %s", webhookEvent.EventType)
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func (h *WebhookHandler) handleVerificationCompleted(ctx context.Context, event models.WebhookEvent) {
	log.Printf("Verification completed for session: %s", event.Data.SessionID)

	// Extract user_id from metadata or user_data
	userID := h.extractUserID(event.Data)
	if userID == "" {
		log.Printf("User ID not found in webhook data")
		return
	}

	sessions, err := h.repo.ListVerificationSessionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting sessions for user: %v", err)
		return
	}

	if len(sessions) == 0 {
		log.Printf("No sessions found for user: %s", userID)
		return
	}

	// Use the most recent session.
	session := sessions[0]

	// Update using the didit_session_id
	_, err = h.repo.UpdateDiditSessionID(ctx, session.SessionID, event.Data.SessionID)
	if err != nil {
		log.Printf("Error updating didit session data: %v", err)
		return
	}

	emailUser := service.User{
		Name:  session.UserFirstName,
		Email: session.UserEmail,
	}
	h.email.SendApprovedKycEmailAsync(emailUser)
	_, err = h.repo.UpdateStatus(ctx, session.SessionID, "approved")
	if err != nil {
		log.Printf("Error updating session status: %v", err)
		return
	}
	log.Printf("User %s verification approved (session: %s)", session.UserID, session.SessionID)
}

func (h *WebhookHandler) handleVerificationFailed(ctx context.Context, event models.WebhookEvent) {
	log.Printf("Verification failed for session: %s", event.Data.SessionID)

	userID := h.extractUserID(event.Data)
	if userID == "" {
		log.Printf("User ID not found in webhook data")
		return
	}

	sessions, err := h.repo.ListVerificationSessionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting sessions for user: %v", err)
		return
	}

	if len(sessions) == 0 {
		log.Printf("No sessions found for user: %s", userID)
		return
	}

	session := sessions[0]

	emailUser := service.User{
		Name:  session.UserFirstName,
		Email: session.UserEmail,
	}
	h.email.SendFailedKycEmailAsync(emailUser)
	_, err = h.repo.UpdateDiditSessionID(ctx, session.SessionID, event.Data.SessionID)
	if err != nil {
		log.Printf("Error updating didit session data: %v", err)
		return
	}

	_, err = h.repo.UpdateStatus(ctx, session.SessionID, "failed")
	if err != nil {
		log.Printf("Error updating session status: %v", err)
		return
	}

	log.Printf("User %s verification failed (session: %s)", session.UserID, session.SessionID)
}

func (h *WebhookHandler) handleVerificationReview(ctx context.Context, event models.WebhookEvent) {
	log.Printf("Verification under review for session: %s", event.Data.SessionID)

	userID := h.extractUserID(event.Data)
	if userID == "" {
		log.Printf("User ID not found in webhook data")
		return
	}

	sessions, err := h.repo.ListVerificationSessionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting sessions for user: %v", err)
		return
	}

	if len(sessions) == 0 {
		log.Printf("No sessions found for user: %s", userID)
		return
	}

	session := sessions[0]

	_, err = h.repo.UpdateDiditSessionID(ctx, session.SessionID, event.Data.SessionID)
	if err != nil {
		log.Printf("Error updating didit session data: %v", err)
		return
	}

	_, err = h.repo.UpdateStatus(ctx, session.SessionID, "review")
	if err != nil {
		log.Printf("Error updating session status: %v", err)
		return
	}

	log.Printf("User %s verification under review (session: %s)", session.UserID, session.SessionID)
}

// extractUserID extracts the user_id from the webhook data
func (h *WebhookHandler) extractUserID(data models.WebhookData) string {
	// First, try the user_id field directly.
	if data.UserID != "" {
		return data.UserID
	}
	// Then try the metadata.
	if data.Metadata != nil {
		if userID, ok := data.Metadata["user_id"]; ok {
			if str, ok := userID.(string); ok {
				return str
			}
		}
		if userID, ok := data.Metadata["internal_user_id"]; ok {
			if str, ok := userID.(string); ok {
				return str
			}
		}
	}

	// Finally try the user_data.
	if data.UserData != nil {
		if userID, ok := data.UserData["user_id"]; ok {
			if str, ok := userID.(string); ok {
				return str
			}
		}
		if userID, ok := data.UserData["id"]; ok {
			if str, ok := userID.(string); ok {
				return str
			}
		}
	}

	return ""
}

// StartVerification initiates the KYC verification process using the fixed link.
func (h *WebhookHandler) StartVerification(c *gin.Context) {
	var req models.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := uuid.New().String()

	// Create a session in the database.
	session := &models.VerificationSession{
		UserID:        req.UserID,
		SessionID:     sessionID,
		Status:        "pending",
		UserEmail:     req.Email,
		UserFirstName: req.FirstName,
		UserLastName:  req.LastName,
	}

	_, err := h.repo.CreateSession(c.Request.Context(), session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Generate verification URL using Didit's fixed link.
	verificationURL := h.diditClient.GetVerificationURL(req.UserID, req.Email, req.FirstName, req.LastName)

	response := models.VerificationResponse{
		VerificationURL: verificationURL,
		UserID:          req.UserID,
	}

	log.Printf("Verification session created for user %s: %s", req.UserID, verificationURL)

	c.JSON(http.StatusOK, response)
}

func (h *WebhookHandler) GetVerificationStatus(c *gin.Context) {
	sessionID := c.Param("sessionId")

	session, err := h.repo.GetSessionBySessionID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetUserVerifications retrieves all verifications for a user.
func (h *WebhookHandler) GetUserVerifications(c *gin.Context) {
	userID := c.Param("userId")

	sessions, err := h.repo.ListVerificationSessionsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// GetUserVerificationStatus retrieves the status of a user's last verification.
func (h *WebhookHandler) GetUserVerificationStatus(c *gin.Context) {
	userID := c.Param("userId")

	sessions, err := h.repo.ListVerificationSessionsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(sessions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No verification sessions found for user"})
		return
	}

	// Return to the most recent session
	c.JSON(http.StatusOK, sessions[0])
}
