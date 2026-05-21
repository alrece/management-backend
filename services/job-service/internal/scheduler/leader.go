package scheduler

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const leaderKey = "job:scheduler:leader"

// LeaderElection Redis Leader 选举（SETNX + 续期 + failover）
type LeaderElection struct {
	rdb      *redis.Client
	interval time.Duration
	ttl      time.Duration
	instance string
	logger   *zap.Logger
	isLeader bool
	stopCh   chan struct{}
}

func NewLeaderElection(rdb *redis.Client, logger *zap.Logger) *LeaderElection {
	hostname, _ := os.Hostname()
	return &LeaderElection{
		rdb:      rdb,
		interval: 8 * time.Second,
		ttl:      30 * time.Second,
		instance: fmt.Sprintf("%s-%d", hostname, os.Getpid()),
		logger:   logger,
		stopCh:   make(chan struct{}),
	}
}

func (le *LeaderElection) Start(ctx context.Context, onLeader func()) {
	ticker := time.NewTicker(le.interval)
	defer ticker.Stop()

	// 首次立即尝试
	if le.tryAcquire(ctx) {
		le.logger.Info("获得调度器 Leader", zap.String("instance", le.instance))
		onLeader()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-le.stopCh:
			return
		case <-ticker.C:
			if le.IsLeader() {
				le.renew(ctx)
			} else {
				if le.tryAcquire(ctx) {
					le.logger.Info("获得调度器 Leader", zap.String("instance", le.instance))
					onLeader()
				}
			}
		}
	}
}

func (le *LeaderElection) tryAcquire(ctx context.Context) bool {
	ok, err := le.rdb.SetNX(ctx, leaderKey, le.instance, le.ttl).Result()
	if err != nil {
		le.logger.Error("Leader 选举失败", zap.Error(err))
		return false
	}
	le.isLeader = ok
	return ok
}

func (le *LeaderElection) renew(ctx context.Context) {
	val, err := le.rdb.Get(ctx, leaderKey).Result()
	if err != nil || val != le.instance {
		le.isLeader = false
		return
	}
	le.rdb.Expire(ctx, leaderKey, le.ttl)
}

func (le *LeaderElection) IsLeader() bool {
	return le.isLeader
}

func (le *LeaderElection) Stop() {
	close(le.stopCh)
}
