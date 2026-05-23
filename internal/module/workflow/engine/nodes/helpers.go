package nodes

import "fmt"

// strVal 从 config map 中提取字符串值
func strVal(config map[string]any, key string) string {
	v, ok := config[key]
	if !ok || v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// intVal 从 config map 中提取整数值，支持默认值
func intVal(config map[string]any, key string, defaultVal int) int {
	v, ok := config[key]
	if !ok || v == nil {
		return defaultVal
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	default:
		return defaultVal
	}
}
