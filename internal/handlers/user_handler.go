package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/FelipePn10/crispaybackend/internal/email/service"
)

type EmailHandler struct {
	emailService *service.EmailService
}

func NewEmailHandler(emailService *service.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

type EmailRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *EmailHandler) SendApprovedEmailKYC(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest

	// Decode the JSON from the request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
	}
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
	}

	emailUser := service.User{
		Name:  req.Name,
		Email: req.Email,
	}
	h.emailService.SendApprovedKycEmailAsync(emailUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "KYC approval email sent successfully.",
		"email":   req.Email,
	})
}
