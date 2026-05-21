package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"management-backend/services/workflow-service/internal/config"

	"github.com/sony/gobreaker"
)

// n8n API 请求/响应结构

type CreateWorkflowReq struct {
	Name string `json:"name"`
	// n8n workflow 创建所需的最小字段
	Nodes []map[string]interface{} `json:"nodes,omitempty"`
}

type WorkflowResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Active bool `json:"active"`
}

type ExecuteReq struct {
	// n8n 执行请求参数
	WorkflowData map[string]interface{} `json:"workflowData,omitempty"`
	StartNodes   []map[string]interface{} `json:"startNodes,omitempty"`
	RunData      map[string]interface{} `json:"runData,omitempty"`
}

type ExecutionResp struct {
	ID string `json:"id"`
}

// N8nClient n8n REST API 封装（带 Circuit Breaker）
type N8nClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	cb         *gobreaker.CircuitBreaker
}

// NewN8nClient 创建 n8n 客户端
func NewN8nClient(cfg config.N8NConfig) *N8nClient {
	cbSettings := gobreaker.Settings{
		Name:        "n8n",
		MaxRequests: 3,                              // 半开状态允许的探测请求数
		Interval:    60 * time.Second,                // 统计间隔
		Timeout:     30 * time.Second,                // Open→Half-Open 等待时间
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5      // 连续 5 次失败触发 Open
		},
	}

	crudTimeout := time.Duration(cfg.Timeout) * time.Second
	if crudTimeout == 0 {
		crudTimeout = 30 * time.Second
	}

	return &N8nClient{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{Timeout: crudTimeout},
		cb:      gobreaker.NewCircuitBreaker(cbSettings),
	}
}

func (c *N8nClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	_, err := c.cb.Execute(func() (interface{}, error) {
		var bodyReader io.Reader
		if body != nil {
			data, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("序列化请求体失败: %w", err)
			}
			bodyReader = bytes.NewReader(data)
		}

		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if c.apiKey != "" {
			req.Header.Set("X-N8N-API-KEY", c.apiKey)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("请求 n8n 失败: %w", err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("n8n 服务错误: %d %s", resp.StatusCode, string(data))
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("n8n 客户端错误: %d %s", resp.StatusCode, string(data))
		}

		return data, nil
	})
	return nil, err
}

// CreateWorkflow 在 n8n 创建工作流
func (c *N8nClient) CreateWorkflow(ctx context.Context, name string) (*WorkflowResp, error) {
	data, err := c.doRequest(ctx, http.MethodPost, "/api/v1/workflows", CreateWorkflowReq{Name: name})
	if err != nil {
		return nil, err
	}
	var resp WorkflowResp
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析 n8n 响应失败: %w", err)
	}
	return &resp, nil
}

// GetWorkflow 获取 n8n 工作流
func (c *N8nClient) GetWorkflow(ctx context.Context, n8nID string) (*WorkflowResp, error) {
	data, err := c.doRequest(ctx, http.MethodGet, "/api/v1/workflows/"+n8nID, nil)
	if err != nil {
		return nil, err
	}
	var resp WorkflowResp
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析 n8n 响应失败: %w", err)
	}
	return &resp, nil
}

// DeleteWorkflow 在 n8n 删除工作流
func (c *N8nClient) DeleteWorkflow(ctx context.Context, n8nID string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/api/v1/workflows/"+n8nID, nil)
	return err
}

// ActivateWorkflow 激活 n8n 工作流
func (c *N8nClient) ActivateWorkflow(ctx context.Context, n8nID string) error {
	_, err := c.doRequest(ctx, http.MethodPost, "/api/v1/workflows/"+n8nID+"/activate", nil)
	return err
}

// DeactivateWorkflow 停用 n8n 工作流
func (c *N8nClient) DeactivateWorkflow(ctx context.Context, n8nID string) error {
	_, err := c.doRequest(ctx, http.MethodPost, "/api/v1/workflows/"+n8nID+"/deactivate", nil)
	return err
}

// ExecuteWorkflow 执行 n8n 工作流
func (c *N8nClient) ExecuteWorkflow(ctx context.Context, n8nID string) (*ExecutionResp, error) {
	// 使用执行超时覆盖默认超时
	execTimeout := time.Duration(config.C.N8N.ExecTimeout) * time.Second
	if execTimeout == 0 {
		execTimeout = 120 * time.Second
	}

	origClient := c.httpClient
	c.httpClient = &http.Client{Timeout: execTimeout}
	defer func() { c.httpClient = origClient }()

	data, err := c.doRequest(ctx, http.MethodPost, "/api/v1/workflows/"+n8nID+"/run", nil)
	if err != nil {
		return nil, err
	}
	var resp ExecutionResp
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析 n8n 响应失败: %w", err)
	}
	return &resp, nil
}

// CircuitBreakerState 返回熔断器当前状态
func (c *N8nClient) CircuitBreakerState() string {
	return c.cb.State().String()
}
