package actions

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/d3m0k1d/BanForge/internal/config"
)

func SendEmail(action config.Action) error {
	if !action.Enabled {
		return nil
	}

	if action.SMTPHost == "" {
		return fmt.Errorf("SMTP host is empty")
	}

	if action.Email == "" {
		return fmt.Errorf("recipient email is empty")
	}

	if action.EmailSender == "" {
		return fmt.Errorf("sender email is empty")
	}

	addr := fmt.Sprintf("%s:%d", action.SMTPHost, action.SMTPPort)

	subject := action.EmailSubject
	if subject == "" {
		subject = "BanForge Alert"
	}

	var body strings.Builder
	body.WriteString("From: " + action.EmailSender + "\r\n")
	body.WriteString("To: " + action.Email + "\r\n")
	body.WriteString("Subject: " + subject + "\r\n")
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	body.WriteString("\r\n")
	body.WriteString(action.Body)

	auth := smtp.PlainAuth("", action.SMTPUser, action.SMTPPassword, action.SMTPHost)

	if action.SMTPTLS {
		return sendEmailWithTLS(
			addr,
			auth,
			action.EmailSender,
			[]string{action.Email},
			body.String(),
		)
	}

	return smtp.SendMail(
		addr,
		auth,
		action.EmailSender,
		[]string{action.Email},
		[]byte(body.String()),
	)
}

func sendEmailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("split host port: %w", err)
	}

	tlsConfig := &tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("dial TLS: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer func() {
		_ = c.Close()
	}()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return fmt.Errorf("SMTP server does not support AUTH")
		}
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("authenticate: %w", err)
		}
	}

	if err = c.Mail(from); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return fmt.Errorf("rcpt to: %w", err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close data: %w", err)
	}

	return c.Quit()
}
