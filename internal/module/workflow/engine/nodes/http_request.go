package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"management-backend/internal/module/workflow/engine"
)

// HTTPRequestHandler HTTP 请求节点
type HTTPRequestHandler struct{}

func NewHTTPRequestHandler() *HTTPRequestHandler {
	return &HTTPRequestHandler{}
}

func (h *HTTPRequestHandler) Type() string     { return "http_request" }
func (h *HTTPRequestHandler) Category() string  { return "action" }

func (h *HTTPRequestHandler) Schema() engine.NodeSchema {
	return engine.NodeSchema{
		Type:     "http_request",
		Label:    "HTTP 请求",
		Category: "action",
		Icon:     "globe",
		Fields: []engine.SchemaField{
			{Name: "url", Label: "URL", Type: "string", Required: true},
			{Name: "method", Label: "请求方法", Type: "select", Required: true, Default: "GET",
				Options: []engine.SelectOption{
					{Value: "GET", Label: "GET"},
					{Value: "POST", Label: "POST"},
					{Value: "PUT", Label: "PUT"},
					{Value: "DELETE", Label: "DELETE"},
					{Value: "PATCH", Label: "PATCH"},
				}},
			{Name: "headers", Label: "请求头", Type: "json"},
			{Name: "body", Label: "请求体", Type: "json"},
			{Name: "timeout", Label: "超时(秒)", Type: "number", Default: 30},
		},
	}
}

func (h *HTTPRequestHandler) Execute(ctx context.Context, input *engine.NodeInput) (*engine.NodeOutput, error) {
	url := strVal(input.Config, "url")
	if url == "" {
		return nil, fmt.Errorf("url 不能为空")
	}
	method := strings.ToUpper(strVal(input.Config, "method"))
	if method == "" {
		method = "GET"
	}

	timeout := time.Duration(intVal(input.Config, "timeout", 30)) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 构建请求体
	var bodyReader io.Reader
	if bodyVal, ok := input.Config["body"]; ok && bodyVal != nil {
		b, _ := json.Marshal(bodyVal)
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	if headers, ok := input.Config["headers"]; ok {
		if m, ok := headers.(map[string]any); ok {
			for k, v := range m {
				req.Header.Set(k, fmt.Sprint(v))
			}
		}
	}
	if bodyReader != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应体为 map
	var data map[string]any
	if err := json.Unmarshal(respBody, &data); err != nil {
		// 非 JSON 响应，包装为字符串
		data = map[string]any{
			"statusCode": resp.StatusCode,
			"body":       string(respBody),
		}
	} else {
		data["statusCode"] = resp.StatusCode
	}

	return &engine.NodeOutput{Items: []map[string]any{data}}, nil
}
