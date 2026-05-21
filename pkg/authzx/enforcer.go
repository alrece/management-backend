package authzx

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/redis/go-redis/v9"
)

// Casbin RBAC 模型：sub=用户/角色, dom=租户ID, perm=权限标识
const rbacModelText = `
[request_definition]
r = sub, dom, perm

[policy_definition]
p = sub, dom, perm

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.perm == p.perm
`

const notifyChannel = "casbin:notify"

// enforcerEntry 缓存条目，记录最后访问时间
type enforcerEntry struct {
	enforcer   *casbin.Enforcer
	lastAccess time.Time
}

// EnforcerManager 多租户 Casbin Enforcer 管理器
type EnforcerManager struct {
	rdb    *redis.Client
	cache  sync.Map // tenantID(int64) -> *enforcerEntry
	stopCh chan struct{}
}

// NewEnforcerManager 创建管理器
func NewEnforcerManager(rdb *redis.Client) *EnforcerManager {
	return &EnforcerManager{
		rdb:    rdb,
		stopCh: make(chan struct{}),
	}
}

// GetEnforcer 获取租户专属 Enforcer（带缓存 + 访问时间更新）
func (m *EnforcerManager) GetEnforcer(tenantID int64) (*casbin.Enforcer, error) {
	if v, ok := m.cache.Load(tenantID); ok {
		entry := v.(*enforcerEntry)
		entry.lastAccess = time.Now()
		return entry.enforcer, nil
	}

	md, err := model.NewModelFromString(rbacModelText)
	if err != nil {
		return nil, fmt.Errorf("创建 Casbin model 失败: %w", err)
	}

	adapter := NewRedisAdapter(m.rdb, tenantID)
	enforcer, err := casbin.NewEnforcer(md, adapter)
	if err != nil {
		return nil, fmt.Errorf("创建 Enforcer 失败: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("加载策略失败: %w", err)
	}

	newEntry := &enforcerEntry{enforcer: enforcer, lastAccess: time.Now()}
	actual, loaded := m.cache.LoadOrStore(tenantID, newEntry)
	if loaded {
		entry := actual.(*enforcerEntry)
		entry.lastAccess = time.Now()
		return entry.enforcer, nil
	}
	return enforcer, nil
}

// Invalidate 清除租户缓存（下次请求重新加载）
func (m *EnforcerManager) Invalidate(tenantID int64) {
	m.cache.Delete(tenantID)
}

// StartEviction 启动周期性淘汰，每隔 interval 清理 idleTimeout 未使用的 enforcer
func (m *EnforcerManager) StartEviction(interval, idleTimeout time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-m.stopCh:
				ticker.Stop()
				return
			case <-ticker.C:
				now := time.Now()
				m.cache.Range(func(key, value any) bool {
					entry := value.(*enforcerEntry)
					if now.Sub(entry.lastAccess) > idleTimeout {
						m.cache.Delete(key)
					}
					return true
				})
			}
		}
	}()
}

// Stop 停止淘汰协程
func (m *EnforcerManager) Stop() {
	close(m.stopCh)
}

// CheckPermission 检查用户权限
func (m *EnforcerManager) CheckPermission(ctx context.Context, userID, tenantID int64, permission string) (bool, error) {
	enforcer, err := m.GetEnforcer(tenantID)
	if err != nil {
		return false, err
	}
	return enforcer.Enforce(
		strconv.FormatInt(userID, 10),
		strconv.FormatInt(tenantID, 10),
		permission,
	)
}

// SyncRolePermissions 同步角色权限到 Casbin
func (m *EnforcerManager) SyncRolePermissions(ctx context.Context, tenantID int64, roleCode string, permissions []string) error {
	enforcer, err := m.GetEnforcer(tenantID)
	if err != nil {
		return err
	}

	dom := strconv.FormatInt(tenantID, 10)
	enforcer.RemoveFilteredPolicy(0, roleCode, dom)

	for _, perm := range permissions {
		enforcer.AddPolicy(roleCode, dom, perm)
	}

	return m.Notify(ctx, tenantID)
}

// SyncUserRoles 同步用户角色到 Casbin
func (m *EnforcerManager) SyncUserRoles(ctx context.Context, tenantID, userID int64, roleCodes []string) error {
	enforcer, err := m.GetEnforcer(tenantID)
	if err != nil {
		return err
	}

	userSub := strconv.FormatInt(userID, 10)
	dom := strconv.FormatInt(tenantID, 10)

	enforcer.RemoveFilteredGroupingPolicy(0, userSub, "", dom)
	for _, code := range roleCodes {
		enforcer.AddGroupingPolicy(userSub, code, dom)
	}

	return m.Notify(ctx, tenantID)
}

// Notify 发布策略变更通知
func (m *EnforcerManager) Notify(ctx context.Context, tenantID int64) error {
	return m.rdb.Publish(ctx, notifyChannel, tenantID).Err()
}
