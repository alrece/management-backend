package nodes

import (
	"context"
	"encoding/json"
	"fmt"

	"management-backend/internal/module/workflow/engine"
	"gorm.io/gorm"
)

// SQLQueryHandler SQL 查询节点（只读查询，禁止写操作）
type SQLQueryHandler struct {
	db *gorm.DB
}

func NewSQLQueryHandler(db *gorm.DB) *SQLQueryHandler {
	return &SQLQueryHandler{db: db}
}

func (h *SQLQueryHandler) Type() string     { return "sql_query" }
func (h *SQLQueryHandler) Category() string  { return "action" }

func (h *SQLQueryHandler) Schema() engine.NodeSchema {
	return engine.NodeSchema{
		Type:     "sql_query",
		Label:    "SQL 查询",
		Category: "action",
		Icon:     "database",
		Fields: []engine.SchemaField{
			{Name: "sql", Label: "SQL 语句", Type: "string", Required: true},
			{Name: "params", Label: "参数", Type: "json"},
			{Name: "maxRows", Label: "最大行数", Type: "number", Default: 100},
		},
	}
}

func (h *SQLQueryHandler) Execute(ctx context.Context, input *engine.NodeInput) (*engine.NodeOutput, error) {
	sql := strVal(input.Config, "sql")
	if sql == "" {
		return nil, fmt.Errorf("sql 不能为空")
	}

	maxRows := intVal(input.Config, "maxRows", 100)
	if maxRows <= 0 || maxRows > 1000 {
		maxRows = 100
	}

	// 提取参数
	var args []any
	if params, ok := input.Config["params"]; ok {
		switch p := params.(type) {
		case []any:
			args = p
		case map[string]any:
			args = make([]any, 0, len(p))
			for _, v := range p {
				args = append(args, v)
			}
		}
	}

	db := h.db.WithContext(ctx).Raw(sql, args...)
	rows, err := db.Rows()
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %w", err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var results []map[string]any
	count := 0

	for rows.Next() && count < maxRows {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, fmt.Errorf("扫描行失败: %w", err)
		}

		row := make(map[string]any, len(cols))
		for i, col := range cols {
			val := values[i]
			// GORM 返回 []byte，转为字符串
			if b, ok := val.([]byte); ok {
				var parsed any
				if json.Unmarshal(b, &parsed) == nil {
					row[col] = parsed
				} else {
					row[col] = string(b)
				}
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
		count++
	}

	if results == nil {
		results = []map[string]any{}
	}

	return &engine.NodeOutput{Items: results}, nil
}
