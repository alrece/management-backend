package nodes

import (
	"context"
	"fmt"
	"strings"

	"management-backend/internal/module/workflow/engine"
)

// TransformHandler 数据转换节点
// 支持 map（字段映射）、filter（条件过滤）、merge（合并）三种模式
type TransformHandler struct{}

func NewTransformHandler() *TransformHandler {
	return &TransformHandler{}
}

func (h *TransformHandler) Type() string     { return "transform" }
func (h *TransformHandler) Category() string  { return "transform" }

func (h *TransformHandler) Schema() engine.NodeSchema {
	return engine.NodeSchema{
		Type:     "transform",
		Label:    "数据转换",
		Category: "transform",
		Icon:     "shuffle",
		Fields: []engine.SchemaField{
			{Name: "mode", Label: "转换模式", Type: "select", Required: true, Default: "map",
				Options: []engine.SelectOption{
					{Value: "map", Label: "字段映射"},
					{Value: "filter", Label: "条件过滤"},
					{Value: "jsonPath", Label: "JSON 路径提取"},
				}},
			{Name: "mappings", Label: "映射规则", Type: "json"},
			{Name: "filterExpr", Label: "过滤表达式", Type: "string"},
			{Name: "path", Label: "JSON 路径", Type: "string"},
		},
	}
}

func (h *TransformHandler) Execute(ctx context.Context, input *engine.NodeInput) (*engine.NodeOutput, error) {
	mode := strVal(input.Config, "mode")
	if mode == "" {
		mode = "map"
	}

	switch mode {
	case "map":
		return h.doMap(input)
	case "filter":
		return h.doFilter(input)
	case "jsonPath":
		return h.doJsonPath(input)
	default:
		return nil, fmt.Errorf("不支持的转换模式: %s", mode)
	}
}

// doMap 字段映射转换
func (h *TransformHandler) doMap(input *engine.NodeInput) (*engine.NodeOutput, error) {
	mappings, ok := input.Config["mappings"]
	if !ok {
		return nil, fmt.Errorf("map 模式需要 mappings 配置")
	}

	mappingMap, ok := mappings.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("mappings 必须是 key-value 对象")
	}

	var results []map[string]any
	for _, item := range input.Items {
		row := make(map[string]any, len(mappingMap))
		for targetField, sourceField := range mappingMap {
			key := fmt.Sprint(sourceField)
			row[targetField] = getNestedField(item, key)
		}
		results = append(results, row)
	}

	if results == nil {
		results = []map[string]any{}
	}
	return &engine.NodeOutput{Items: results}, nil
}

// doFilter 条件过滤
func (h *TransformHandler) doFilter(input *engine.NodeInput) (*engine.NodeOutput, error) {
	filterExpr := strVal(input.Config, "filterExpr")
	if filterExpr == "" {
		return nil, fmt.Errorf("filter 模式需要 filterExpr 配置")
	}

	// 简单字段匹配过滤：field=value 格式
	parts := strings.SplitN(filterExpr, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("filterExpr 格式错误，应为 field=value")
	}
	filterField := strings.TrimSpace(parts[0])
	filterValue := strings.TrimSpace(parts[1])

	var results []map[string]any
	for _, item := range input.Items {
		val := fmt.Sprint(getNestedField(item, filterField))
		if val == filterValue {
			results = append(results, item)
		}
	}

	if results == nil {
		results = []map[string]any{}
	}
	return &engine.NodeOutput{Items: results}, nil
}

// doJsonPath JSON 路径提取
func (h *TransformHandler) doJsonPath(input *engine.NodeInput) (*engine.NodeOutput, error) {
	path := strVal(input.Config, "path")
	if path == "" {
		return nil, fmt.Errorf("jsonPath 模式需要 path 配置")
	}

	var results []map[string]any
	for _, item := range input.Items {
		val := getNestedField(item, path)
		// 将提取的值包装为单字段结果
		results = append(results, map[string]any{"value": val})
	}

	if results == nil {
		results = []map[string]any{}
	}
	return &engine.NodeOutput{Items: results}, nil
}

