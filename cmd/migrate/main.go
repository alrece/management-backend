package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"management-backend/internal/config"
	"management-backend/internal/multitenant"

	"go.uber.org/zap"
)

func main() {
	cfgPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	target := flag.Int64("target", 0, "指定租户ID迁移（0=迁移默认库）")
	all := flag.Bool("all", false, "迁移所有已就绪租户")
	flag.Parse()

	if err := config.Load(*cfgPath); err != nil {
		fmt.Fprintf(os.Stderr, "配置加载失败: %v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	migrator := multitenant.NewMigrator()
	ctx := context.Background()

	switch {
	case *all:
		logger.Info("开始迁移所有租户")
		if err := migrator.MigrateAll(ctx); err != nil {
			logger.Fatal("迁移失败", zap.Error(err))
		}
		logger.Info("全部迁移完成")
	case *target > 0:
		logger.Info("迁移租户", zap.Int64("tenantID", *target))
		if err := migrator.MigrateTenant(ctx, *target); err != nil {
			logger.Fatal("迁移失败", zap.Error(err))
		}
		logger.Info("租户迁移完成")
	default:
		logger.Info("迁移默认库")
		if err := migrator.MigrateDefault(ctx); err != nil {
			logger.Fatal("迁移失败", zap.Error(err))
		}
		logger.Info("默认库迁移完成")
	}
}
