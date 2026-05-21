package authzx

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// PolicySync Redis Pub/Sub 策略同步器
type PolicySync struct {
	rdb     *redis.Client
	manager *EnforcerManager
	sub     *redis.PubSub
}

// NewPolicySync 创建同步器
func NewPolicySync(rdb *redis.Client, manager *EnforcerManager) *PolicySync {
	return &PolicySync{rdb: rdb, manager: manager}
}

// Start 启动后台订阅
func (s *PolicySync) Start(ctx context.Context) error {
	s.sub = s.rdb.Subscribe(ctx, notifyChannel)

	go func() {
		ch := s.sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				if tenantID, err := strconv.ParseInt(msg.Payload, 10, 64); err == nil {
					s.manager.Invalidate(tenantID)
				}
			}
		}
	}()

	return nil
}

// Stop 停止订阅
func (s *PolicySync) Stop() error {
	if s.sub != nil {
		return s.sub.Close()
	}
	return nil
}
