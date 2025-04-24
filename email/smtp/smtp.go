package smtp

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/core-stack/authz/email"
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

	templates email.TemplateRegistry
}

func New(cfg SMTPConfig, templates email.TemplateRegistry) *SMTPSender {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &SMTPSender{cfg: cfg, auth: auth, templates: templates}
}

func (s *SMTPSender) Send(ctx context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	// Headers do e-mail
	headers := map[string]string{
		"From":         s.cfg.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"utf-8\"",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)
	return smtp.SendMail(addr, s.auth, s.cfg.From, []string{to}, []byte(msg.String()))
}

func (s *SMTPSender) SendActiveAccount(ctx context.Context, user zmodel.User, code string) error {
	return s.sendWithTemplateCode(ctx, email.TemplateActiveAccount, user, code)
}

func (s *SMTPSender) SendResetPassword(ctx context.Context, user zmodel.User, code string) error {
	return s.sendWithTemplateCode(ctx, email.TemplateResetPassword, user, code)
}

func (s *SMTPSender) SendNotifyResetPassword(ctx context.Context, user zmodel.User) error {
	return s.sendWithTemplate(ctx, email.TemplateActiveAccount, user)
}

func (s *SMTPSender) SendChangePassword(ctx context.Context, user zmodel.User) error {
	return s.sendWithTemplate(ctx, email.TemplateResetPassword, user)
}

func (s *SMTPSender) SendDeleteAccount(ctx context.Context, user zmodel.User) error {
	return s.sendWithTemplate(ctx, email.TemplateActiveAccount, user)
}

func (s *SMTPSender) sendWithTemplateCode(
	ctx context.Context,
	templateType email.TemplateType,
	user zmodel.User,
	code string,
) error {
	tpl, ok := s.templates[templateType]
	if !ok {
		return fmt.Errorf("template for %s not found", templateType)
	}

	data := struct {
		Code string
		Name string
	}{
		Name: user.Name,
		Code: code,
	}

	var body bytes.Buffer
	if err := tpl.Template.Execute(&body, data); err != nil {
		return err
	}
	return s.Send(ctx, user.Email, tpl.Subject, body.String())
}

func (s *SMTPSender) sendWithTemplate(
	ctx context.Context,
	templateType email.TemplateType,
	user zmodel.User,
) error {
	tpl, ok := s.templates[templateType]
	if !ok {
		return fmt.Errorf("template for %s not found", templateType)
	}

	data := struct {
		Name string
	}{
		Name: user.Name,
	}

	var body bytes.Buffer
	if err := tpl.Template.Execute(&body, data); err != nil {
		return err
	}
	return s.Send(ctx, user.Email, tpl.Subject, body.String())
}
