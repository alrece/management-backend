package authzx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/casbin/casbin/v2/model"
	"github.com/redis/go-redis/v9"
)

// RedisAdapter Casbin Redis 存储适配器（每租户独立）
type RedisAdapter struct {
	rdb      *redis.Client
	tenantID int64
}

// NewRedisAdapter 创建 Redis 适配器
func NewRedisAdapter(rdb *redis.Client, tenantID int64) *RedisAdapter {
	return &RedisAdapter{rdb: rdb, tenantID: tenantID}
}

func (a *RedisAdapter) policyKey(ptype string) string {
	return fmt.Sprintf("casbin:%s:%d", ptype, a.tenantID)
}

// LoadPolicy 从 Redis 加载策略到 Casbin model
func (a *RedisAdapter) LoadPolicy(md model.Model) error {
	ctx := context.Background()

	pRules, err := a.loadRules(ctx, "p")
	if err != nil {
		return err
	}
	for _, rule := range pRules {
		md.AddPolicy("p", "p", rule)
	}

	gRules, err := a.loadRules(ctx, "g")
	if err != nil {
		return err
	}
	for _, rule := range gRules {
		md.AddPolicy("g", "g", rule)
	}

	return nil
}

// SavePolicy 将 Casbin model 策略保存到 Redis
func (a *RedisAdapter) SavePolicy(md model.Model) error {
	ctx := context.Background()

	pPolicies, _ := md.GetPolicy("p", "p")
	if err := a.saveRules(ctx, "p", pPolicies); err != nil {
		return err
	}
	gPolicies, _ := md.GetPolicy("g", "g")
	return a.saveRules(ctx, "g", gPolicies)
}

// AddPolicy 添加单条策略
func (a *RedisAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	ctx := context.Background()
	rules, _ := a.loadRules(ctx, ptype)
	rules = append(rules, rule)
	return a.saveRules(ctx, ptype, rules)
}

// RemovePolicy 删除单条策略
func (a *RedisAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	ctx := context.Background()
	rules, _ := a.loadRules(ctx, ptype)
	filtered := make([][]string, 0, len(rules))
	for _, r := range rules {
		if !equalRule(r, rule) {
			filtered = append(filtered, r)
		}
	}
	return a.saveRules(ctx, ptype, filtered)
}

// RemoveFilteredPolicy 按条件删除策略
func (a *RedisAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	ctx := context.Background()
	rules, _ := a.loadRules(ctx, ptype)
	filtered := make([][]string, 0, len(rules))
	for _, r := range rules {
		if !matchFilter(r, fieldIndex, fieldValues) {
			filtered = append(filtered, r)
		}
	}
	return a.saveRules(ctx, ptype, filtered)
}

func (a *RedisAdapter) loadRules(ctx context.Context, ptype string) ([][]string, error) {
	data, err := a.rdb.Get(ctx, a.policyKey(ptype)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var rules [][]string
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func (a *RedisAdapter) saveRules(ctx context.Context, ptype string, rules [][]string) error {
	data, err := json.Marshal(rules)
	if err != nil {
		return err
	}
	return a.rdb.Set(ctx, a.policyKey(ptype), data, 0).Err()
}

func equalRule(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func matchFilter(rule []string, fieldIndex int, fieldValues []string) bool {
	for i, v := range fieldValues {
		if v == "" {
			continue
		}
		idx := fieldIndex + i
		if idx >= len(rule) || rule[idx] != v {
			return false
		}
	}
	return true
}
