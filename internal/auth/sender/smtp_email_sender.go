package sender

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromName    string
	FromAddress string
}

type SMTPEmailSender struct {
	config SMTPConfig
}

func NewSMTPEmailSender(config SMTPConfig) (*SMTPEmailSender, error) {
	config.Host = strings.TrimSpace(config.Host)
	config.Username = strings.TrimSpace(config.Username)
	config.FromName = strings.TrimSpace(config.FromName)
	config.FromAddress = strings.TrimSpace(config.FromAddress)

	if config.Host == "" {
		return nil, fmt.Errorf("smtp host is required")
	}
	if config.Port == 0 {
		config.Port = 587
	}
	if config.Port < 1 || config.Port > 65535 {
		return nil, fmt.Errorf("smtp port is invalid: %d", config.Port)
	}
	if config.FromAddress == "" {
		return nil, fmt.Errorf("smtp from address is required")
	}
	if _, err := mail.ParseAddress(config.FromAddress); err != nil {
		return nil, fmt.Errorf("smtp from address is invalid: %w", err)
	}

	return &SMTPEmailSender{config: config}, nil
}

func (s *SMTPEmailSender) Send(to, subject, content string) error {
	to = strings.TrimSpace(to)
	if to == "" {
		return fmt.Errorf("smtp recipient address is required")
	}
	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("smtp recipient address is invalid: %w", err)
	}

	message, err := s.buildMessage(to, subject, content)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port))
	if s.config.Port == 465 {
		return s.sendWithTLS(addr, to, message)
	}

	if err := smtp.SendMail(addr, s.auth(), s.config.FromAddress, []string{to}, message); err != nil {
		return fmt.Errorf("smtp send mail failed: %w", err)
	}
	return nil
}

func (s *SMTPEmailSender) auth() smtp.Auth {
	if s.config.Username == "" {
		return nil
	}
	return smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
}

func (s *SMTPEmailSender) buildMessage(to, subject, content string) ([]byte, error) {
	var body bytes.Buffer
	qp := quotedprintable.NewWriter(&body)
	if _, err := qp.Write([]byte(content)); err != nil {
		return nil, fmt.Errorf("encode smtp body failed: %w", err)
	}
	if err := qp.Close(); err != nil {
		return nil, fmt.Errorf("finalize smtp body failed: %w", err)
	}

	fromHeader := s.config.FromAddress
	if s.config.FromName != "" {
		fromHeader = (&mail.Address{
			Name:    s.config.FromName,
			Address: s.config.FromAddress,
		}).String()
	}

	headers := []string{
		fmt.Sprintf("From: %s", fromHeader),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", mime.QEncoding.Encode("UTF-8", subject)),
		fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: quoted-printable",
		"",
		body.String(),
	}

	return []byte(strings.Join(headers, "\r\n")), nil
}

func (s *SMTPEmailSender) sendWithTLS(addr, to string, message []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: s.config.Host,
		MinVersion: tls.VersionTLS12,
	})
	if err != nil {
		return fmt.Errorf("smtp tls dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("smtp client init failed: %w", err)
	}
	defer client.Close()

	if auth := s.auth(); auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
	}
	if err := client.Mail(s.config.FromAddress); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO failed: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return fmt.Errorf("smtp write message failed: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("smtp finalize message failed: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("smtp quit failed: %w", err)
	}
	return nil
}
