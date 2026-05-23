package nodes

import (
	"context"
	"fmt"
	"strings"

	"management-backend/internal/module/workflow/engine"
)

// ConditionHandler 条件分支节点
// 根据 condition 表达式判断走 true 还是 false 分支
// 由执行引擎根据 sourceHandle（true/false）路由到下游节点
type ConditionHandler struct{}

func NewConditionHandler() *ConditionHandler {
	return &ConditionHandler{}
}

func (h *ConditionHandler) Type() string     { return "condition" }
func (h *ConditionHandler) Category() string  { return "control" }

func (h *ConditionHandler) Schema() engine.NodeSchema {
	return engine.NodeSchema{
		Type:     "condition",
		Label:    "条件判断",
		Category: "control",
		Icon:     "git-branch",
		Fields: []engine.SchemaField{
			{Name: "field", Label: "字段路径", Type: "string", Required: true},
			{Name: "operator", Label: "比较运算符", Type: "select", Required: true, Default: "eq",
				Options: []engine.SelectOption{
					{Value: "eq", Label: "等于"},
					{Value: "neq", Label: "不等于"},
					{Value: "gt", Label: "大于"},
					{Value: "gte", Label: "大于等于"},
					{Value: "lt", Label: "小于"},
					{Value: "lte", Label: "小于等于"},
					{Value: "contains", Label: "包含"},
					{Value: "empty", Label: "为空"},
					{Value: "notEmpty", Label: "不为空"},
				}},
			{Name: "value", Label: "比较值", Type: "string"},
		},
	}
}

func (h *ConditionHandler) Execute(ctx context.Context, input *engine.NodeInput) (*engine.NodeOutput, error) {
	field := strVal(input.Config, "field")
	operator := strVal(input.Config, "operator")
	compareVal := strVal(input.Config, "value")

	if field == "" {
		return nil, fmt.Errorf("field 不能为空")
	}

	// 从上游输入数据中提取字段值
	var actualVal any
	if len(input.Items) > 0 {
		actualVal = getNestedField(input.Items[0], field)
	}

	result := evaluateCondition(actualVal, operator, compareVal)

	// 输出标记 branch，执行引擎据此路由
	branch := "false"
	if result {
		branch = "true"
	}

	return &engine.NodeOutput{
		Items: []map[string]any{
			{"branch": branch, "field": field, "matched": result},
		},
	}, nil
}

// evaluateCondition 执行条件判断
func evaluateCondition(actual any, operator, expected string) bool {
	actualStr := fmt.Sprint(actual)

	switch operator {
	case "eq":
		return actualStr == expected
	case "neq":
		return actualStr != expected
	case "gt":
		return actualStr > expected
	case "gte":
		return actualStr >= expected
	case "lt":
		return actualStr < expected
	case "lte":
		return actualStr <= expected
	case "contains":
		return strings.Contains(actualStr, expected)
	case "empty":
		return actual == nil || actualStr == "" || actualStr == "<nil>"
	case "notEmpty":
		return actual != nil && actualStr != "" && actualStr != "<nil>"
	default:
		return false
	}
}

// getNestedField 支持点号路径访问嵌套字段（如 "data.name"）
func getNestedField(m map[string]any, path string) any {
	parts := strings.Split(path, ".")
	var cur any = m
	for _, p := range parts {
		switch v := cur.(type) {
		case map[string]any:
			cur = v[p]
		default:
			return nil
		}
	}
	return cur
}
