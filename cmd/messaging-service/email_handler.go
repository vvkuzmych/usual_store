package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

//go:embed templates
var emailTemplatesFS embed.FS

// EmailHandler implements the messaging.EmailHandler interface
type EmailHandler struct {
	smtpConfig smtpConfig
	infoLog    *log.Logger
	errorLog   *log.Logger
}

type smtpConfig struct {
	host     string
	port     int
	username string
	password string
}

// NewEmailHandler creates a new EmailHandler
func NewEmailHandler(smtp smtpConfig, infoLog, errorLog *log.Logger) *EmailHandler {
	return &EmailHandler{
		smtpConfig: smtp,
		infoLog:    infoLog,
		errorLog:   errorLog,
	}
}

// SendEmail sends an email using SMTP
func (h *EmailHandler) SendEmail(from, to, subject, tmpl string, data map[string]interface{}) error {
	// Render HTML template
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)
	t, err := template.New("email-html").ParseFS(emailTemplatesFS, templateToRender)
	if err != nil {
		h.errorLog.Printf("Error parsing HTML template: %v", err)
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var htmlBuffer bytes.Buffer
	if err = t.ExecuteTemplate(&htmlBuffer, "body", data); err != nil {
		h.errorLog.Printf("Error executing HTML template: %v", err)
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}
	htmlMessage := htmlBuffer.String()

	// Render plain text template
	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplatesFS, templateToRender)
	if err != nil {
		h.errorLog.Printf("Error parsing plain template: %v", err)
		return fmt.Errorf("failed to parse plain template: %w", err)
	}

	var plainBuffer bytes.Buffer
	if err = t.ExecuteTemplate(&plainBuffer, "body", data); err != nil {
		h.errorLog.Printf("Error executing plain template: %v", err)
		return fmt.Errorf("failed to execute plain template: %w", err)
	}
	plainMessage := plainBuffer.String()

	// Configure SMTP server
	server := mail.NewSMTPClient()
	server.Host = h.smtpConfig.host
	server.Port = h.smtpConfig.port
	server.Username = h.smtpConfig.username
	server.Password = h.smtpConfig.password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Connect to SMTP server
	smtpClient, err := server.Connect()
	if err != nil {
		h.errorLog.Printf("Error connecting to SMTP server: %v", err)
		return fmt.Errorf("failed to connect to SMTP: %w", err)
	}

	// Create email message
	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, htmlMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		h.errorLog.Printf("Error sending email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	h.infoLog.Printf("Email sent successfully: From=%s, To=%s, Subject=%s", from, to, subject)
	return nil
}
