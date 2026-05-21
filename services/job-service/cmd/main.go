package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"management-backend/services/job-service/internal/config"
	"management-backend/services/job-service/internal/handler"
	"management-backend/services/job-service/internal/repository"
	"management-backend/services/job-service/internal/router"
	"management-backend/services/job-service/internal/scheduler"
	"management-backend/services/job-service/internal/service"

	"management-backend/pkg/authx"
	"management-backend/pkg/discoveryx"
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
	if err := config.Load("configs/config.yaml"); err != nil {
		fmt.Fprintf(os.Stderr, "配置加载失败: %v\n", err)
		os.Exit(1)
	}

	zapLogger := initZap()
	defer zapLogger.Sync()

	if err := snowflake.Init(config.C.Snowflake.Node); err != nil {
		zapLogger.Fatal("雪花ID初始化失败", zap.Error(err))
	}

	// MySQL
	db, err := gorm.Open(mysql.Open(config.C.MySQL.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		zapLogger.Fatal("MySQL 连接失败", zap.Error(err))
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.Addr(),
		Password: config.C.Redis.Password,
		DB:       config.C.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		zapLogger.Fatal("Redis 连接失败", zap.Error(err))
	}

	// 依赖注入
	taskRepo := repository.NewJobTaskRepo(db)
	logRepo := repository.NewExecLogRepo(db)
	jobSvc := service.NewJobService(taskRepo, logRepo)
	jobHandler := handler.NewJobHandler(jobSvc)

	// Scheduler
	sched := scheduler.NewScheduler(taskRepo, logRepo, rdb, zapLogger)
	leaderElection := scheduler.NewLeaderElection(rdb, zapLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go leaderElection.Start(ctx, func() {
		if err := sched.Start(ctx); err != nil {
			zapLogger.Error("调度器启动失败", zap.Error(err))
		}
	})

	// HTTP 服务
	gin.SetMode(config.C.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Consul 服务注册
	if config.C.Consul.Enabled {
		registry, err := discoveryx.NewConsulRegistry(discoveryx.ConsulConfig{
			Enabled:     true,
			Address:     config.C.Consul.Address,
			ServiceAddr: config.C.Consul.ServiceAddr,
		}, zapLogger)
		if err != nil {
			zapLogger.Fatal("Consul 注册器创建失败", zap.Error(err))
		}
		if err := registry.Register("mb-job-service", config.C.Server.Port); err != nil {
			zapLogger.Fatal("Consul 注册失败", zap.Error(err))
		}
		defer registry.Deregister()
	}

	authed := engine.Group("/api")
	jwtSecret := config.C.JWT.Secret
	authed.Use(func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未授权"})
			return
		}
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}
		claims, err := authx.ParseToken(tokenStr, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "Token 无效"})
			return
		}
		c.Set("userId", claims.UserID)
		c.Set("tenantId", claims.TenantID)
		c.Next()
	})
	router.RegisterRoutes(authed, jobHandler)

	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine}

	go func() {
		zapLogger.Info("job-service 启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("收到关闭信号...")

	sched.Stop()
	leaderElection.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)

	zapLogger.Info("job-service 已关闭")
}

func initZap() *zap.Logger {
	zapCfg := zap.NewProductionConfig()
	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, _ := zapCfg.Build()
	return l
}
