package nodes

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"management-backend/internal/module/workflow/engine"
)

// SendEmailHandler 发送邮件节点
type SendEmailHandler struct {
	// SMTP 配置（从外部注入，不暴露到节点配置）
	host     string
	port     string
	username string
	password string
	from     string
}

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSendEmailHandler(cfg EmailConfig) *SendEmailHandler {
	return &SendEmailHandler{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
	}
}

func (h *SendEmailHandler) Type() string     { return "send_email" }
func (h *SendEmailHandler) Category() string  { return "action" }

func (h *SendEmailHandler) Schema() engine.NodeSchema {
	return engine.NodeSchema{
		Type:     "send_email",
		Label:    "发送邮件",
		Category: "action",
		Icon:     "mail",
		Fields: []engine.SchemaField{
			{Name: "to", Label: "收件人", Type: "string", Required: true},
			{Name: "cc", Label: "抄送", Type: "string"},
			{Name: "subject", Label: "主题", Type: "string", Required: true},
			{Name: "body", Label: "正文", Type: "string", Required: true},
			{Name: "contentType", Label: "内容类型", Type: "select", Default: "text/plain",
				Options: []engine.SelectOption{
					{Value: "text/plain", Label: "纯文本"},
					{Value: "text/html", Label: "HTML"},
				}},
		},
	}
}

func (h *SendEmailHandler) Execute(ctx context.Context, input *engine.NodeInput) (*engine.NodeOutput, error) {
	if h.host == "" {
		return nil, fmt.Errorf("SMTP 未配置")
	}

	to := strVal(input.Config, "to")
	subject := strVal(input.Config, "subject")
	body := strVal(input.Config, "body")
	contentType := strVal(input.Config, "contentType")
	if contentType == "" {
		contentType = "text/plain"
	}

	if to == "" || subject == "" {
		return nil, fmt.Errorf("收件人和主题不能为空")
	}

	toList := strings.Split(to, ",")
	for i := range toList {
		toList[i] = strings.TrimSpace(toList[i])
	}

	// 构建邮件内容
	var msg strings.Builder
	msg.WriteString("From: " + h.from + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	if cc := strVal(input.Config, "cc"); cc != "" {
		msg.WriteString("Cc: " + cc + "\r\n")
	}
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("Content-Type: " + contentType + "; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := h.host + ":" + h.port
	auth := smtp.PlainAuth("", h.username, h.password, h.host)

	if err := smtp.SendMail(addr, auth, h.from, toList, []byte(msg.String())); err != nil {
		return nil, fmt.Errorf("发送邮件失败: %w", err)
	}

	return &engine.NodeOutput{
		Items: []map[string]any{
			{"sent": true, "to": to, "subject": subject},
		},
	}, nil
}
