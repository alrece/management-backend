package engine

import (
	"context"
	"fmt"
	"sync"
)

// NodeHandler 节点处理器接口
type NodeHandler interface {
	// Execute 执行节点逻辑
	Execute(ctx context.Context, input *NodeInput) (*NodeOutput, error)
	// Type 返回节点类型标识
	Type() string
	// Category 返回节点分类（trigger / action / control / transform）
	Category() string
	// Schema 返回节点 Schema（供前端属性面板使用）
	Schema() NodeSchema
}

// NodeSchema 节点 Schema 定义（供前端渲染属性面板和节点面板）
type NodeSchema struct {
	Type     string        `json:"type"`
	Label    string        `json:"label"`
	Category string        `json:"category"`
	Icon     string        `json:"icon"`
	Fields   []SchemaField `json:"fields"`
}

// SchemaField Schema 字段定义
type SchemaField struct {
	Name     string         `json:"name"`
	Label    string         `json:"label"`
	Type     string         `json:"type"` // string, number, boolean, select, json
	Required bool           `json:"required"`
	Default  any            `json:"default,omitempty"`
	Options  []SelectOption `json:"options,omitempty"`
}

// SelectOption 下拉选项
type SelectOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// NodeRegistry 节点注册表
type NodeRegistry struct {
	mu       sync.RWMutex
	handlers map[string]NodeHandler
}

// NewNodeRegistry 创建节点注册表
func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		handlers: make(map[string]NodeHandler),
	}
}

// Register 注册节点处理器
func (r *NodeRegistry) Register(h NodeHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.handlers[h.Type()]; exists {
		panic(fmt.Sprintf("节点类型 %q 已注册", h.Type()))
	}
	r.handlers[h.Type()] = h
}

// Get 获取节点处理器
func (r *NodeRegistry) Get(nodeType string) (NodeHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[nodeType]
	return h, ok
}

// AllSchemas 返回所有已注册节点的 Schema
func (r *NodeRegistry) AllSchemas() []NodeSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]NodeSchema, 0, len(r.handlers))
	for _, h := range r.handlers {
		result = append(result, h.Schema())
	}
	return result
}

// SchemasByCategory 按分类返回 Schema
func (r *NodeRegistry) SchemasByCategory(category string) []NodeSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []NodeSchema
	for _, h := range r.handlers {
		if h.Category() == category {
			result = append(result, h.Schema())
		}
	}
	return result
}
