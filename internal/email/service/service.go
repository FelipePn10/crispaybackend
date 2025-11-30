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
    <title>Bem-vindo ao CrisPay</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #F7F5F3;">
    <table width="100%" cellpadding="0" cellspacing="0" style="background-color: #F7F5F3; padding: 40px 20px;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 9px; box-shadow: 0px 0px 0px 0.9px rgba(0,0,0,0.08); overflow: hidden; max-width: 100%;">
                    
                    <tr>
                        <td style="padding: 42px 40px 32px 40px; border-bottom: 1px solid rgba(55,50,47,0.06);">
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td>
                                        <h1 style="color: #2F3037; margin: 0; font-size: 20px; font-weight: 500; letter-spacing: -0.01em;">
                                            CrisPay
                                        </h1>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <tr>
                        <td style="padding: 40px 40px 24px 40px;">
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center">
                                        <!-- Badge -->
                                        <table cellpadding="0" cellspacing="0" style="margin-bottom: 20px;">
                                            <tr>
                                                <td style="padding: 6px 14px; background-color: #ffffff; border: 1px solid rgba(2,6,23,0.08); border-radius: 90px; box-shadow: 0px 0px 0px 4px rgba(55,50,47,0.05);">
                                                    <span style="color: #37322F; font-size: 12px; font-weight: 500; line-height: 12px;">
                                                        Seu KYC foi aprovado. Parabéns!
                                                    </span>
                                                </td>
                                            </tr>
                                        </table>
                                        
                                        <!-- Título -->
                                        <h2 style="color: #49423D; margin: 0 0 16px 0; font-size: 36px; font-weight: 600; line-height: 1.2; letter-spacing: -0.02em; text-align: center;">
                                            Sua conta está pronta, {{.Name}}
                                        </h2>
                                        
                                        <!-- Descrição -->
                                        <p style="color: #605A57; font-size: 16px; line-height: 28px; margin: 0; text-align: center; max-width: 480px;">
                                            Agora você já pode pagar com cripto em qualquer lugar do mundo. Rápido, seguro e transparente.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <tr>
                        <td style="padding: 0 40px 32px 40px;">
                            <table width="100%" cellpadding="0" cellspacing="0" style="border-top: 1px solid #E0DEDB; border-bottom: 1px solid #E0DEDB;">
                                <!-- Card 1 -->
                                <tr>
                                    <td style="padding: 24px 0; border-bottom: 1px solid rgba(224,222,219,0.5);">
                                        <h3 style="color: #49423D; margin: 0 0 8px 0; font-size: 14px; font-weight: 600; line-height: 24px;">
                                            Total transparência
                                        </h3>
                                        <p style="color: #605A57; font-size: 13px; line-height: 22px; margin: 0;">
                                            Acompanhe cada etapa da sua compra: pagamento, conversão e envio. Receba comprovantes e rastreie seu pedido diretamente da loja.
                                        </p>
                                    </td>
                                </tr>
                                
                                <tr>
                                    <td style="padding: 24px 0; border-bottom: 1px solid rgba(224,222,219,0.5);">
                                        <h3 style="color: #49423D; margin: 0 0 8px 0; font-size: 14px; font-weight: 600; line-height: 24px;">
                                            Compre com cripto, sem conversão necessária
                                        </h3>
                                        <p style="color: #605A57; font-size: 13px; line-height: 22px; margin: 0;">
                                            Não precisa mais converter BTC, ETH ou USDT para gastar. Pague em cripto — nós cuidamos do resto.
                                        </p>
                                    </td>
                                </tr>
                                
                                <!-- Card 3 -->
                                <tr>
                                    <td style="padding: 24px 0;">
                                        <h3 style="color: #49423D; margin: 0 0 8px 0; font-size: 14px; font-weight: 600; line-height: 24px;">
                                            Simplicidade que inspira confiança
                                        </h3>
                                        <p style="color: #605A57; font-size: 13px; line-height: 22px; margin: 0;">
                                            Taxas claras, processo simplificado e suporte humano. Uma experiência de pagamento cripto projetada para quem valoriza conveniência.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <tr>
                        <td align="center" style="padding: 32px 40px 40px 40px;">
                            <table cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center" style="background-color: #37322F; border-radius: 50px; box-shadow: 0px 1px 2px rgba(55,50,47,0.12);">
                                        <a href="#" style="background-color: #37322F; color: #ffffff; padding: 12px 32px; text-decoration: none; font-size: 14px; font-weight: 500; display: inline-block; border-radius: 50px;">
                                            Começar a usar
                                        </a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <tr>
                        <td style="padding: 0 40px 40px 40px;">
                            <p style="color: #828387; font-size: 14px; line-height: 24px; margin: 0; text-align: center;">
                                Precisa de ajuda? Nossa equipe está sempre disponível para você.
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #F7F5F3; padding: 32px 40px; border-top: 1px solid rgba(55,50,47,0.06);">
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center">
                                        <p style="color: #828387; font-size: 13px; margin: 0 0 8px 0;">
                                            © 2025 CrisPay. Todos os direitos reservados.
                                        </p>
                                        <p style="color: #828387; font-size: 12px; margin: 0;">
                                            Você recebeu este email porque teve seu KYC aprovado - CrisPay.
                                        </p>
                                    </td>
                                </tr>
                                
                                <!-- Links do footer -->
                                <tr>
                                    <td align="center" style="padding-top: 20px;">
                                        <table cellpadding="0" cellspacing="0">
                                            <tr>
                                                <td style="padding: 0 12px;">
                                                    <a href="#" style="color: #605A57; font-size: 13px; text-decoration: none;">
                                                        Produtos
                                                    </a>
                                                </td>
                                                <td style="padding: 0 12px; border-left: 1px solid rgba(55,50,47,0.12);">
                                                    <a href="#" style="color: #605A57; font-size: 13px; text-decoration: none;">
                                                        Pricing
                                                    </a>
                                                </td>
                                                <td style="padding: 0 12px; border-left: 1px solid rgba(55,50,47,0.12);">
                                                    <a href="#" style="color: #605A57; font-size: 13px; text-decoration: none;">
                                                        Docs
                                                    </a>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>
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
    <title>Atualização sobre seu KYC - CrisPay</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #F7F5F3;">
    <table width="100%" cellpadding="0" cellspacing="0" style="background-color: #F7F5F3; padding: 40px 20px;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 9px; box-shadow: 0px 0px 0px 0.9px rgba(0,0,0,0.08); overflow: hidden; max-width: 100%;">
                    
                    <!-- Header -->
                    <tr>
                        <td style="padding: 42px 40px 32px 40px; border-bottom: 1px solid rgba(55,50,47,0.06);">
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td>
                                        <h1 style="color: #2F3037; margin: 0; font-size: 20px; font-weight: 500; letter-spacing: -0.01em;">
                                            CrisPay
                                        </h1>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <!-- Hero Section -->
                    <tr>
                        <td style="padding: 40px 40px 24px 40px;">
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center">
                                        <!-- Badge de Atenção -->
                                        <table cellpadding="0" cellspacing="0" style="margin-bottom: 20px;">
                                            <tr>
                                                <td style="padding: 6px 14px; background-color: #FFF8F0; border: 1px solid rgba(217, 119, 6, 0.2); border-radius: 90px; box-shadow: 0px 0px 0px 4px rgba(217, 119, 6, 0.05);">
                                                    <span style="color: #92400E; font-size: 12px; font-weight: 500; line-height: 12px;">
                                                        Ação necessária
                                                    </span>
                                                </td>
                                            </tr>
                                        </table>
                                        
                                        <!-- Título -->
                                        <h2 style="color: #49423D; margin: 0 0 16px 0; font-size: 36px; font-weight: 600; line-height: 1.2; letter-spacing: -0.02em; text-align: center;">
                                            Não conseguimos verificar sua identidade
                                        </h2>
                                        
                                        <!-- Descrição -->
                                        <p style="color: #605A57; font-size: 16px; line-height: 28px; margin: 0; text-align: center; max-width: 480px;">
                                            Olá, {{.Name}}. Infelizmente, não foi possível concluir a verificação da sua conta. Veja os detalhes abaixo.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <!-- Motivos da Rejeição -->
                    <tr>
                        <td style="padding: 0 40px 32px 40px;">
                            <table width="100%" cellpadding="0" cellspacing="0" style="border-top: 1px solid #E0DEDB; border-bottom: 1px solid #E0DEDB;">
                                <!-- Card 1 -->
                                <tr>
                                    <td style="padding: 24px 0; border-bottom: 1px solid rgba(224,222,219,0.5);">
                                        <h3 style="color: #49423D; margin: 0 0 8px 0; font-size: 14px; font-weight: 600; line-height: 24px;">
                                            Motivo da rejeição
                                        </h3>
                                        <p style="color: #605A57; font-size: 13px; line-height: 22px; margin: 0;">
                                            {{.RejectionReason}}
                                        </p>
                                    </td>
                                </tr>
                                
                                <!-- Card 2 -->
                                <tr>
                                    <td style="padding: 24px 0; border-bottom: 1px solid rgba(224,222,219,0.5);">
                                        <h3 style="color: #49423D; margin: 0 0 8px 0; font-size: 14px; font-weight: 600; line-height: 24px;">
                                            O que você precisa fazer
                                        </h3>
                                        <p style="color: #605A57; font-size: 13px; line-height: 22px; margin: 0;">
                                            Envie novos documentos com melhor qualidade: foto nítida, boa iluminação, e certifique-se de que todos os dados estão legíveis.
                                        </p>
                                    </td>
                                </tr>
                                
                                <!-- Card 3 -->
                                <tr>
                                    <td style="padding: 24px 0;">
                                        <h3 style="color: #49423D; margin: 0 0 8px 0; font-size: 14px; font-weight: 600; line-height: 24px;">
                                            Precisa de ajuda?
                                        </h3>
                                        <p style="color: #605A57; font-size: 13px; line-height: 22px; margin: 0;">
                                            Nossa equipe de suporte está disponível para esclarecer dúvidas e ajudá-lo a completar o processo com sucesso.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <!-- CTAs -->
                    <tr>
                        <td align="center" style="padding: 32px 40px 40px 40px;">
                            <table cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center" style="background-color: #37322F; border-radius: 50px; box-shadow: 0px 1px 2px rgba(55,50,47,0.12); margin-bottom: 16px;">
                                        <a href="#" style="background-color: #37322F; color: #ffffff; padding: 12px 32px; text-decoration: none; font-size: 14px; font-weight: 500; display: inline-block; border-radius: 50px;">
                                            Enviar novos documentos
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <table cellpadding="0" cellspacing="0" style="margin-top: 12px;">
                                <tr>
                                    <td align="center">
                                        <a href="#" style="color: #605A57; padding: 8px 16px; text-decoration: none; font-size: 14px; font-weight: 500; display: inline-block;">
                                            Falar com o suporte →
                                        </a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                    <!-- Mensagem de suporte -->
                    <tr>
                        <td style="padding: 0 40px 40px 40px;">
                            <p style="color: #828387; font-size: 14px; line-height: 24px; margin: 0; text-align: center;">
                                Estamos aqui para ajudar você a completar este processo.<br/>Responda este email ou acesse nossa central de ajuda.
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #F7F5F3; padding: 32px 40px; border-top: 1px solid rgba(55,50,47,0.06);">
                            <table width="100%" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center">
                                        <p style="color: #828387; font-size: 13px; margin: 0 0 8px 0;">
                                            © 2025 CrisPay. Todos os direitos reservados.
                                        </p>
                                        <p style="color: #828387; font-size: 12px; margin: 0;">
                                            Você recebeu este email sobre a verificação da sua conta CrisPay.
                                        </p>
                                    </td>
                                </tr>
                                
                                <!-- Links do footer -->
                                <tr>
                                    <td align="center" style="padding-top: 20px;">
                                        <table cellpadding="0" cellspacing="0">
                                            <tr>
                                                <td style="padding: 0 12px;">
                                                    <a href="#" style="color: #605A57; font-size: 13px; text-decoration: none;">
                                                        Produtos
                                                    </a>
                                                </td>
                                                <td style="padding: 0 12px; border-left: 1px solid rgba(55,50,47,0.12);">
                                                    <a href="#" style="color: #605A57; font-size: 13px; text-decoration: none;">
                                                        Pricing
                                                    </a>
                                                </td>
                                                <td style="padding: 0 12px; border-left: 1px solid rgba(55,50,47,0.12);">
                                                    <a href="#" style="color: #605A57; font-size: 13px; text-decoration: none;">
                                                        Docs
                                                    </a>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
	tmpl, err := template.New("reproved").Parse(htmlTemplate)
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
