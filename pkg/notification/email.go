package notification

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"github.com/jagac/pfinance/pkg/config"
)

// EmailNotifier struct holds the configuration required to send email notifications.
// It uses the configuration from the GlobalConfig struct (like SMTP server settings).
type EmailNotifier struct {
	config *config.GlobalConfig // Configuration for email settings (e.g., SMTP server, credentials)
}

// NewEmailNotifier creates a new instance of EmailNotifier with the provided config.
func NewEmailNotifier(config *config.GlobalConfig) *EmailNotifier {
	return &EmailNotifier{
		config: config,
	}
}

// Send sends an email using SMTP, authenticating via TLS and SMTP credentials.
// It accepts a context to allow cancellation, recipient email, subject, and body of the email.
func (e *EmailNotifier) Send(ctx context.Context, to, subject, body string) error {

	select {
	case <-ctx.Done():
		return nil
	default:

	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", e.config.From, to, subject, body)

	auth := smtp.PlainAuth("", e.config.From, e.config.Password, e.config.SMTPHost)

	conn, err := smtp.Dial(fmt.Sprintf("%s:%s", e.config.SMTPHost, e.config.SMTPPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.config.SMTPHost,
	}

	if err := conn.StartTLS(tlsConfig); err != nil {
		return err
	}

	if err := conn.Auth(auth); err != nil {
		return err
	}

	if err := conn.Mail(e.config.From); err != nil {
		return err
	}

	if err := conn.Rcpt(to); err != nil {
		return err
	}

	wc, err := conn.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return err
	}

	return wc.Close()
}
