package nodes

import (
	"context"
	"fmt"
	"time"

	"management-backend/internal/module/workflow/engine"
)

// DelayHandler 延迟节点
type DelayHandler struct{}

func NewDelayHandler() *DelayHandler {
	return &DelayHandler{}
}

func (h *DelayHandler) Type() string     { return "delay" }
func (h *DelayHandler) Category() string  { return "control" }

func (h *DelayHandler) Schema() engine.NodeSchema {
	return engine.NodeSchema{
		Type:     "delay",
		Label:    "延时等待",
		Category: "control",
		Icon:     "clock",
		Fields: []engine.SchemaField{
			{Name: "duration", Label: "等待时间(秒)", Type: "number", Required: true, Default: 1},
		},
	}
}

func (h *DelayHandler) Execute(ctx context.Context, input *engine.NodeInput) (*engine.NodeOutput, error) {
	seconds := intVal(input.Config, "duration", 1)
	if seconds <= 0 {
		seconds = 1
	}
	if seconds > 300 {
		return nil, fmt.Errorf("延时不能超过 300 秒，当前: %d", seconds)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(seconds) * time.Second):
		// 透传上游数据
		items := input.Items
		if items == nil {
			items = []map[string]any{}
		}
		return &engine.NodeOutput{
			Items: items,
		}, nil
	}
}
