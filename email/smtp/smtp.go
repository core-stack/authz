package smtp

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/core-stack/authz/zmodel"
)

type SMTPConfig struct {
	Host     string // exemplo: "smtp.gmail.com"
	Port     int    // exemplo: 587
	Username string
	Password string
	From     string // exemplo: "noreply@dominio.com"
}

type SMTPSender struct {
	cfg  SMTPConfig
	auth smtp.Auth
}

func New(cfg SMTPConfig) *SMTPSender {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &SMTPSender{cfg: cfg, auth: auth}
}

func (s *SMTPSender) Send(ctx context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	// Headers do e-mail
	headers := map[string]string{
		"From":         s.cfg.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=\"utf-8\"",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)
	return smtp.SendMail(addr, s.auth, s.cfg.From, []string{to}, []byte(msg.String()))
}

func (s *SMTPSender) SendActiveAccount(ctx context.Context, user zmodel.User, code string) error {
	t, err := template.ParseFiles("")
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	// diretorio
	// data

}
