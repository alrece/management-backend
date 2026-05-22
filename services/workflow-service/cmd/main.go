package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"management-backend/services/workflow-service/internal/client"
	"management-backend/services/workflow-service/internal/config"
	"management-backend/services/workflow-service/internal/handler"
	"management-backend/services/workflow-service/internal/repository"
	"management-backend/services/workflow-service/internal/router"
	"management-backend/services/workflow-service/internal/service"

	pb "management-backend/api/proto/workflow"
	"management-backend/pkg/authx"
	"management-backend/pkg/discoveryx"
	"management-backend/pkg/snowflake"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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

	// n8n Client
	n8nClient := client.NewN8nClient(config.C.N8N)

	// 依赖注入
	catRepo := repository.NewCategoryRepo(db)
	wfRepo := repository.NewWorkflowRepo(db)
	instRepo := repository.NewInstanceRepo(db)

	catSvc := service.NewCategoryService(catRepo)
	wfSvc := service.NewWorkflowService(wfRepo, instRepo, n8nClient)

	wfHandler := handler.NewWorkflowHandler(wfSvc, catSvc)
	webhookHandler := handler.NewWebhookHandler(instRepo, wfRepo)

	// HTTP 服务
	gin.SetMode(config.C.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Webhook 路由（无需认证）
	engine.POST("/api/workflow/callback", webhookHandler.Callback)

	// 认证路由
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
	router.RegisterRoutes(authed, wfHandler, webhookHandler)

	// gRPC 服务（端口 9083）
	grpcSrv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 600}))
	pb.RegisterWorkflowServiceServer(grpcSrv, handler.NewWorkflowGRPCHandler(wfSvc, instRepo))

	grpcAddr := fmt.Sprintf(":%d", 9083)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		zapLogger.Fatal("gRPC 监听失败", zap.Error(err))
	}
	go func() {
		zapLogger.Info("workflow-service gRPC 启动", zap.String("addr", grpcAddr))
		if err := grpcSrv.Serve(grpcListener); err != nil {
			zapLogger.Error("gRPC 服务异常", zap.Error(err))
		}
	}()

	// Consul 注册
	if config.C.Consul.Enabled {
		registry, err := discoveryx.NewConsulRegistry(discoveryx.ConsulConfig{
			Enabled:     true,
			Address:     config.C.Consul.Address,
			ServiceAddr: config.C.Consul.ServiceAddr,
		}, zapLogger)
		if err != nil {
			zapLogger.Fatal("Consul 注册器创建失败", zap.Error(err))
		}
		if err := registry.Register("mb-workflow-service", config.C.Server.Port); err != nil {
			zapLogger.Fatal("Consul 注册失败", zap.Error(err))
		}
		defer registry.Deregister()
	}

	// 启动时对账
	go func() {
		time.Sleep(5 * time.Second) // 等待 n8n 就绪
		if err := wfSvc.StartupReconcile(context.Background()); err != nil {
			zapLogger.Warn("启动对账失败", zap.Error(err))
		}
	}()

	// 定期同步（5min）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go startPeriodicSync(ctx, wfSvc, zapLogger)

	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine}

	go func() {
		zapLogger.Info("workflow-service 启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLogger.Info("收到关闭信号...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	grpcSrv.GracefulStop()

	zapLogger.Info("workflow-service 已关闭")
}

func startPeriodicSync(ctx context.Context, wfSvc service.WorkflowService, logger *zap.Logger) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := wfSvc.SyncFromN8n(ctx); err != nil {
				logger.Warn("n8n 同步失败", zap.Error(err))
			}
		}
	}
}

func initZap() *zap.Logger {
	zapCfg := zap.NewProductionConfig()
	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, _ := zapCfg.Build()
	return l
}
