package grpc

import (
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientPool gRPC 连接池（指数退避重连）
type ClientPool struct {
	mu     sync.RWMutex
	conns  map[string]*grpc.ClientConn
	target string
	opts   []grpc.DialOption
}

// NewClientPool 创建连接池
func NewClientPool(target string) *ClientPool {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	}
	return &ClientPool{
		conns:  make(map[string]*grpc.ClientConn),
		target: target,
		opts:   opts,
	}
}

// Get 获取连接（带指数退避重连）
func (p *ClientPool) Get(serviceName string) (*grpc.ClientConn, error) {
	p.mu.RLock()
	if conn, ok := p.conns[serviceName]; ok {
		p.mu.RUnlock()
		return conn, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	if conn, ok := p.conns[serviceName]; ok {
		return conn, nil
	}

	addr := fmt.Sprintf("%s/%s", p.target, serviceName)
	var conn *grpc.ClientConn
	var err error

	for attempt := 0; attempt < 3; attempt++ {
		conn, err = grpc.NewClient(addr, p.opts...)
		if err == nil {
			p.conns[serviceName] = conn
			return conn, nil
		}
		time.Sleep(time.Duration(1<<attempt) * time.Second)
	}
	return nil, fmt.Errorf("连接 %s 失败（重试3次）: %w", addr, err)
}

// Close 关闭所有连接
func (p *ClientPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for name, conn := range p.conns {
		conn.Close()
		delete(p.conns, name)
	}
}
