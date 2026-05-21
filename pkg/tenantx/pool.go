package tenantx

import (
	"container/list"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type poolEntry struct {
	tenantID   int64
	lastAccess time.Time
}

// TenantPool 租户连接池 LRU 管理
type TenantPool struct {
	mu          sync.Mutex
	maxSize     int
	idleTimeout time.Duration
	lru         *list.List
	items       map[int64]*list.Element

	poolSize       atomic.Int64
	evictionsTotal atomic.Int64
}

// NewTenantPool 创建租户连接池
func NewTenantPool(maxSize int, idleTimeout time.Duration) *TenantPool {
	return &TenantPool{
		maxSize:     maxSize,
		idleTimeout: idleTimeout,
		lru:         list.New(),
		items:       make(map[int64]*list.Element),
	}
}

// Touch 更新租户访问时间
func (p *TenantPool) Touch(tenantID int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, ok := p.items[tenantID]; ok {
		entry := elem.Value.(*poolEntry)
		entry.lastAccess = time.Now()
		p.lru.MoveToFront(elem)
	} else {
		entry := &poolEntry{tenantID: tenantID, lastAccess: time.Now()}
		elem := p.lru.PushFront(entry)
		p.items[tenantID] = elem
	}
	p.poolSize.Store(int64(len(p.items)))
}

// EvictOldest 驱逐最久未使用的租户
func (p *TenantPool) EvictOldest() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	elem := p.lru.Back()
	if elem == nil {
		return 0
	}

	entry := elem.Value.(*poolEntry)
	p.lru.Remove(elem)
	delete(p.items, entry.tenantID)
	p.poolSize.Store(int64(len(p.items)))
	p.evictionsTotal.Add(1)
	return entry.tenantID
}

func (p *TenantPool) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.cleanup()
		}
	}
}

func (p *TenantPool) cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	var toEvict []int64

	for elem := p.lru.Back(); elem != nil; elem = elem.Prev() {
		entry := elem.Value.(*poolEntry)
		if now.Sub(entry.lastAccess) > p.idleTimeout {
			toEvict = append(toEvict, entry.tenantID)
		} else {
			break
		}
	}

	for _, id := range toEvict {
		if elem, ok := p.items[id]; ok {
			p.lru.Remove(elem)
			delete(p.items, id)
			p.evictionsTotal.Add(1)
		}
	}
	p.poolSize.Store(int64(len(p.items)))
}

// PoolMetrics 连接池指标
type PoolMetrics struct {
	Size      int64 `json:"size"`
	MaxSize   int   `json:"maxSize"`
	Evictions int64 `json:"evictionsTotal"`
}

// GetMetrics 获取连接池指标
func (p *TenantPool) GetMetrics() PoolMetrics {
	return PoolMetrics{
		Size:      p.poolSize.Load(),
		MaxSize:   p.maxSize,
		Evictions: p.evictionsTotal.Load(),
	}
}
