package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"management-backend/internal/config"
	"management-backend/internal/router"
	"management-backend/pkg/middleware"
	"management-backend/pkg/snowflake"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	if err := config.Load("configs/config.yaml"); err != nil {
		fmt.Fprintf(os.Stderr, "配置加载失败: %v\n", err)
		os.Exit(1)
	}

	// 敏感配置校验（开发模式可跳过）
	if config.C.Server.Mode != "debug" {
		if err := config.C.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "配置校验失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 初始化 Zap 日志
	zapLogger := initLogger(config.C.Log)
	defer zapLogger.Sync()

	// 初始化雪花 ID
	if err := snowflake.Init(config.C.Snowflake.Node); err != nil {
		zapLogger.Fatal("雪花ID初始化失败", zap.Error(err))
	}

	// 连接 MySQL
	db, err := gorm.Open(mysql.Open(config.C.MySQL.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		zapLogger.Fatal("MySQL 连接失败", zap.Error(err))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(config.C.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.C.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 连接 Redis（支持 Sentinel）
	var rdb *redis.Client
	if config.C.Redis.Sentinel.Enabled {
		rdb = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    config.C.Redis.Sentinel.MasterName,
			SentinelAddrs: config.C.Redis.Sentinel.Addrs,
			SentinelPassword: config.C.Redis.Sentinel.Password,
			Password:      config.C.Redis.Password,
			DB:            config.C.Redis.DB,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.C.Redis.Addr(),
			Password: config.C.Redis.Password,
			DB:       config.C.Redis.DB,
		})
	}
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		zapLogger.Fatal("Redis 连接失败", zap.Error(err))
	}

	// 创建引擎
	gin.SetMode(config.C.Server.Mode)
	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recovery(zapLogger))

	// 注册路由
	cleanup := router.Setup(engine, db, rdb, zapLogger)

	// 启动服务
	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine}

	go func() {
		zapLogger.Info("system-service 启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("收到关闭信号，开始优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("服务关闭失败", zap.Error(err))
	}

	sqlDB.Close()
	rdb.Close()
	cleanup()
	zapLogger.Info("服务已关闭")
}

func initLogger(cfg config.LogConfig) *zap.Logger {
	var zapCfg zap.Config
	if cfg.Level == "debug" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
		zapCfg.Encoding = "json"
	}
	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapCfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	l, _ := zapCfg.Build()
	return l
}
