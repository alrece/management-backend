package nodes

import (
	"context"
	"log/slog"

	"management-backend/internal/module/workflow/engine"
)

// LogMessageHandler 日志消息节点（用于调试和工作流标记）
type LogMessageHandler struct{}

func NewLogMessageHandler() *LogMessageHandler {
	return &LogMessageHandler{}
}

func (h *LogMessageHandler) Type() string     { return "log_message" }
func (h *LogMessageHandler) Category() string  { return "action" }

func (h *LogMessageHandler) Schema() engine.NodeSchema {
	return engine.NodeSchema{
		Type:     "log_message",
		Label:    "日志消息",
		Category: "action",
		Icon:     "file-text",
		Fields: []engine.SchemaField{
			{Name: "message", Label: "消息内容", Type: "string", Required: true},
			{Name: "level", Label: "日志级别", Type: "select", Default: "info",
				Options: []engine.SelectOption{
					{Value: "debug", Label: "DEBUG"},
					{Value: "info", Label: "INFO"},
					{Value: "warn", Label: "WARN"},
					{Value: "error", Label: "ERROR"},
				}},
		},
	}
}

func (h *LogMessageHandler) Execute(ctx context.Context, input *engine.NodeInput) (*engine.NodeOutput, error) {
	message := strVal(input.Config, "message")
	level := strVal(input.Config, "level")
	if level == "" {
		level = "info"
	}

	// 记录结构化日志
	logger := slog.With(
		"nodeId", input.NodeID,
		"nodeType", input.NodeType,
	)

	switch level {
	case "debug":
		logger.Debug(message)
	case "warn":
		logger.Warn(message)
	case "error":
		logger.Error(message)
	default:
		logger.Info(message)
	}

	// 透传上游数据
	items := input.Items
	if items == nil {
		items = []map[string]any{}
	}

	return &engine.NodeOutput{
		Items: append(items, map[string]any{
			"loggedMessage": message,
			"level":         level,
		}),
	}, nil
}
