package didit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/FelipePn10/crispaybackend/config"
)

type Client struct {
	config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

// VerifyWebhookSignature validates the webhook signature using HMAC-SHA256
func (c *Client) VerifyWebhookSignature(payload []byte, signature string) bool {
	if signature == "" || c.config.DiditWebhookSecret == "" {
		return false
	}

	h := hmac.New(sha256.New, []byte(c.config.DiditWebhookSecret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare signatures (protection against timing attacks)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// GetVerificationURL returns the fixed URL of the Didit workflow.
func (c *Client) GetVerificationURL(userID, email, firstName, lastName string) string {
	baseURL := "https://verify.didit.me/verify/ynzf6TUsQ6BPDyVMWe8btQ"

	return fmt.Sprintf("%s?user_id=%s&email=%s&first_name=%s&last_name=%s",
		baseURL, userID, email, firstName, lastName)
}
