package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// GracefulShutdown 监听系统信号，执行优雅关闭回调
func GracefulShutdown(logger *zap.Logger, timeout time.Duration, callbacks ...func(ctx context.Context) error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("收到关闭信号", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for i, cb := range callbacks {
		if err := cb(ctx); err != nil {
			logger.Error("关闭回调失败", zap.Int("step", i+1), zap.Error(err))
		}
	}

	logger.Info("服务已优雅关闭")
}
