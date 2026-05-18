package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
)

type snowflakeGenerator struct {
	node *sf.Node
}

// NewGenerator 创建雪花 ID 生成器（返回接口，便于测试 mock）
func NewGenerator(nodeID int64) (IDGenerator, error) {
	node, err := sf.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &snowflakeGenerator{node: node}, nil
}

func (g *snowflakeGenerator) NextID() (int64, error) {
	return g.node.Generate().Int64(), nil
}

// 全局实例
var global IDGenerator

func Init(nodeID int64) error {
	gen, err := NewGenerator(nodeID)
	if err != nil {
		return err
	}
	global = gen
	return nil
}

func NextID() int64 {
	id, _ := global.NextID()
	return id
}
