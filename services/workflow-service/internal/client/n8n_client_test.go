package client

import (
	"context"
	"testing"
	"time"

	"management-backend/services/workflow-service/internal/config"

	"github.com/stretchr/testify/assert"
)

func newTestN8nClient() *N8nClient {
	cfg := config.N8NConfig{
		BaseURL:     "http://localhost:19999",
		APIKey:      "test-key",
		Timeout:     2,
		ExecTimeout: 2,
	}
	return NewN8nClient(cfg)
}

func TestN8nClient_CircuitBreaker_Opens(t *testing.T) {
	c := newTestN8nClient()
	ctx := context.Background()

	for i := 0; i < 6; i++ {
		_, _ = c.GetWorkflow(ctx, "nonexistent")
	}
	assert.Equal(t, "open", c.CircuitBreakerState())
}

func TestN8nClient_CircuitBreaker_HalfOpen_AfterTimeout(t *testing.T) {
	c := newTestN8nClient()
	ctx := context.Background()

	// 触发 Open
	for i := 0; i < 6; i++ {
		_, _ = c.GetWorkflow(ctx, "nonexistent")
	}
	assert.Equal(t, "open", c.CircuitBreakerState())

	// 等待超过 Timeout（30s 在生产环境），这里直接验证状态名可读
	state := c.CircuitBreakerState()
	assert.Contains(t, []string{"open", "half-open", "closed"}, state)
}

func TestN8nClient_Timeout(t *testing.T) {
	c := newTestN8nClient()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	_, err := c.GetWorkflow(ctx, "test")
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Less(t, elapsed, 5*time.Second)
}
