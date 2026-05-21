package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Server gRPC 服务端
type Server struct {
	server   *grpc.Server
	port     int
	listener net.Listener
	logger   *zap.Logger
}

// NewServer 创建 gRPC 服务端（默认 :9081）
func NewServer(port int, logger *zap.Logger, unaryInterceptors []grpc.UnaryServerInterceptor) *Server {
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     300,
			MaxConnectionAge:      600,
			MaxConnectionAgeGrace: 30,
		}),
	}
	if len(unaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	}

	return &Server{
		server: grpc.NewServer(opts...),
		port:   port,
		logger: logger,
	}
}

// GRPCServer 获取底层 grpc.Server 用于注册服务
func (s *Server) GRPCServer() *grpc.Server {
	return s.server
}

// Start 启动 gRPC 服务
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("gRPC 监听失败: %w", err)
	}
	s.listener = listener

	go func() {
		s.logger.Info("gRPC 服务启动", zap.String("addr", addr))
		if err := s.server.Serve(listener); err != nil {
			s.logger.Error("gRPC 服务异常", zap.Error(err))
		}
	}()
	return nil
}

// Stop 优雅停止 gRPC 服务
func (s *Server) Stop() {
	s.server.GracefulStop()
	s.logger.Info("gRPC 服务已停止")
}
