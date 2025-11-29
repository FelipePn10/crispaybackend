package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

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
	htmlTemplate := `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>KYC Aprovado!</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f7fa;">
    <table width="100%" cellpadding="0" cellspacing="0" style="background-color: #f4f7fa; padding: 40px 0;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); overflow: hidden;">
                    <!-- Header com gradiente -->
                    <tr>
                        <td style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 50px 40px; text-align: center;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 32px; font-weight: 600;">
                                🎉 Parábens! Seu KYC acaba de ser aprovado!
                            </h1>
                        </td>
                    </tr>
                    
                    <!-- Conteúdo -->
                    <tr>
                        <td style="padding: 40px;">
                            <h2 style="color: #333333; margin: 0 0 20px 0; font-size: 24px;">
                                Olá, {{.Name}}! 👋
                            </h2>
                            
                            <p style="color: #666666; font-size: 16px; line-height: 1.6; margin: 0 0 20px 0;">
                                Estamos muito felizes em ter você conosco! Sua conta foi verificada com sucesso e você já pode começar a explorar todos os nossos recursos.
                            </p>
                            
                            <p style="color: #666666; font-size: 16px; line-height: 1.6; margin: 0 0 30px 0;">
                                Para começar, aqui estão algumas coisas que você pode fazer:
                            </p>
                            
                            <!-- Cards de recursos -->
                            <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom: 30px;">
                                <tr>
                                    <td style="padding: 20px; background-color: #f8f9fa; border-radius: 8px; margin-bottom: 10px;">
                                        <h3 style="color: #667eea; margin: 0 0 10px 0; font-size: 18px;">
                                            ✨ Faça sua primeira compra
                                        </h3>
                                        <p style="color: #666666; font-size: 14px; line-height: 1.5; margin: 0;">
                                            Conecte sua Wallet para débitos automáticos!
                                        </p>
                                    </td>
                                </tr>
                            </table>
                            
                            <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom: 30px;">
                                <tr>
                                    <td style="padding: 20px; background-color: #f8f9fa; border-radius: 8px;">
                                        <h3 style="color: #667eea; margin: 0 0 10px 0; font-size: 18px;">
                                            🚀 Explore os recursos
                                        </h3>
                                        <p style="color: #666666; font-size: 14px; line-height: 1.5; margin: 0;">
                                            Descubra tudo que nossa plataforma pode fazer por você
                                        </p>
                                    </td>
                                </tr>
                            </table>
                            
                            <!-- Botão CTA -->
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center" style="padding: 20px 0;">
                                        <a href="#" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; padding: 16px 40px; text-decoration: none; border-radius: 8px; font-size: 16px; font-weight: 600; display: inline-block; box-shadow: 0 4px 6px rgba(102, 126, 234, 0.3);">
                                            Acessar Minha Conta
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="color: #999999; font-size: 14px; line-height: 1.6; margin: 30px 0 0 0; text-align: center;">
                                Se você tiver alguma dúvida, não hesite em nos contatar. Estamos aqui para ajudar!
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #f8f9fa; padding: 30px 40px; text-align: center; border-top: 1px solid #e9ecef;">
                            <p style="color: #999999; font-size: 14px; margin: 0 0 10px 0;">
                                © 2025 CrisPay. Todos os direitos reservados.
                            </p>
                            <p style="color: #999999; font-size: 12px; margin: 0;">
                                Você recebeu este email porque concluiu com sucesso o KYC. Não responda a este e-mail.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
	tmpl, err := template.New("approved").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("error parse template")
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, user); err != nil {
		return fmt.Errorf("error execute template")
	}

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.SenderEmail)
	headers["To"] = user.Email
	headers["Subject"] = "KYC Aprovado. Parábens!"
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

func (s *EmailService) SendFailedKycEmail(user User) error {
	htmlTemplate := `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>KYC Aprovado!</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f7fa;">
    <table width="100%" cellpadding="0" cellspacing="0" style="background-color: #f4f7fa; padding: 40px 0;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); overflow: hidden;">
                    <!-- Header com gradiente -->
                    <tr>
                        <td style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 50px 40px; text-align: center;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 32px; font-weight: 600;">
                                🎉 Parábens! Seu KYC acaba de ser aprovado!
                            </h1>
                        </td>
                    </tr>
                    
                    <!-- Conteúdo -->
                    <tr>
                        <td style="padding: 40px;">
                            <h2 style="color: #333333; margin: 0 0 20px 0; font-size: 24px;">
                                Olá, {{.Name}}! 👋
                            </h2>
                            
                            <p style="color: #666666; font-size: 16px; line-height: 1.6; margin: 0 0 20px 0;">
                                Estamos muito felizes em ter você conosco! Sua conta foi verificada com sucesso e você já pode começar a explorar todos os nossos recursos.
                            </p>
                            
                            <p style="color: #666666; font-size: 16px; line-height: 1.6; margin: 0 0 30px 0;">
                                Para começar, aqui estão algumas coisas que você pode fazer:
                            </p>
                            
                            <!-- Cards de recursos -->
                            <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom: 30px;">
                                <tr>
                                    <td style="padding: 20px; background-color: #f8f9fa; border-radius: 8px; margin-bottom: 10px;">
                                        <h3 style="color: #667eea; margin: 0 0 10px 0; font-size: 18px;">
                                            ✨ Faça sua primeira compra
                                        </h3>
                                        <p style="color: #666666; font-size: 14px; line-height: 1.5; margin: 0;">
                                            Conecte sua Wallet para débitos automáticos!
                                        </p>
                                    </td>
                                </tr>
                            </table>
                            
                            <table width="100%" cellpadding="0" cellspacing="0" style="margin-bottom: 30px;">
                                <tr>
                                    <td style="padding: 20px; background-color: #f8f9fa; border-radius: 8px;">
                                        <h3 style="color: #667eea; margin: 0 0 10px 0; font-size: 18px;">
                                            🚀 Explore os recursos
                                        </h3>
                                        <p style="color: #666666; font-size: 14px; line-height: 1.5; margin: 0;">
                                            Descubra tudo que nossa plataforma pode fazer por você
                                        </p>
                                    </td>
                                </tr>
                            </table>
                            
                            <!-- Botão CTA -->
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center" style="padding: 20px 0;">
                                        <a href="#" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; padding: 16px 40px; text-decoration: none; border-radius: 8px; font-size: 16px; font-weight: 600; display: inline-block; box-shadow: 0 4px 6px rgba(102, 126, 234, 0.3);">
                                            Acessar Minha Conta
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="color: #999999; font-size: 14px; line-height: 1.6; margin: 30px 0 0 0; text-align: center;">
                                Se você tiver alguma dúvida, não hesite em nos contatar. Estamos aqui para ajudar!
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #f8f9fa; padding: 30px 40px; text-align: center; border-top: 1px solid #e9ecef;">
                            <p style="color: #999999; font-size: 14px; margin: 0 0 10px 0;">
                                © 2025 CrisPay. Todos os direitos reservados.
                            </p>
                            <p style="color: #999999; font-size: 12px; margin: 0;">
                                Você recebeu este email porque concluiu com sucesso o KYC. Não responda a este e-mail.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
	tmpl, err := template.New("approved").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("error parse template")
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, user); err != nil {
		return fmt.Errorf("error execute template")
	}

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.SenderEmail)
	headers["To"] = user.Email
	headers["Subject"] = "KYC Aprovado. Parábens!"
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
