package service

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

//go:embed templates/*.gohtml
var templatesFS embed.FS

type EmailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderPass  string
	SenderName  string
}

type User struct {
	Name  string
	Email string
}

type EmailService struct {
	config EmailConfig
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{config: config}
}

// Email confirming success with KYC approval.
func (s *EmailService) SendApprovedKycEmail(user User) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/approved_kyc.gohtml")
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, user); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.SenderEmail),
		"To":           user.Email,
		"Subject":      "KYC Aprovado. Parabéns!",
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.Write(body.Bytes())

	auth := smtp.PlainAuth("", s.config.SenderEmail, s.config.SenderPass, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	if err := smtp.SendMail(addr, auth, s.config.SenderEmail, []string{user.Email}, message.Bytes()); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

func (s *EmailService) SendFailedKycEmail(user User) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/failed_kyc.gohtml")
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, user); err != nil {
		return fmt.Errorf("error execute template")
	}

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.SenderEmail)
	headers["To"] = user.Email
	headers["Subject"] = "KYC Reprovado."
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.Write(body.Bytes())

	// Autentication
	auth := smtp.PlainAuth("", s.config.SenderEmail, s.config.SenderPass, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	err = smtp.SendMail(addr, auth, s.config.SenderEmail, []string{user.Email}, message.Bytes())
	if err != nil {
		return fmt.Errorf("error send email: %v", err)
	}
	return nil
}

// SendWelcomeEmailAsync sends the email asynchronously (does not block)
func (s *EmailService) SendApprovedKycEmailAsync(user User) {
	go func() {
		log.Printf("Trying to send an email to: %s", user.Email)
		log.Printf("SMTP Config - Host: %s, Port: %s, Sender: %s",
			s.config.SMTPHost, s.config.SMTPPort, s.config.SenderEmail)

		if err := s.SendApprovedKycEmail(user); err != nil {
			log.Printf("Error sending email to %s: %v", user.Email, err)
		} else {
			log.Printf("Email sent successfully to: %s", user.Email)
		}
	}()
}

func (s *EmailService) SendFailedKycEmailAsync(user User) {
	go func() {
		log.Printf("Trying to send an email to: %s", user.Email)
		log.Printf("SMTP Config - Host: %s, Port: %s, Sender: %s",
			s.config.SMTPHost, s.config.SMTPPort, s.config.SenderEmail)

		if err := s.SendFailedKycEmail(user); err != nil {
			log.Printf("Error sending email to %s: %v", user.Email, err)
		} else {
			log.Printf("Email sent successfully to: %s", user.Email)
		}
	}()
}
