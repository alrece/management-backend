package snowflake

// IDGenerator 雪花 ID 生成器接口（便于测试注入 mock）
type IDGenerator interface {
	NextID() (int64, error)
}
