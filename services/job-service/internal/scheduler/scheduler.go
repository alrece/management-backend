package scheduler

import (
	"context"
	"fmt"
	"time"

	"management-backend/services/job-service/internal/model"
	"management-backend/services/job-service/internal/repository"

	"management-backend/pkg/snowflake"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron     *cron.Cron
	taskRepo repository.JobTaskRepo
	logRepo  repository.ExecLogRepo
	rdb      *redis.Client
	logger   *zap.Logger
}

func NewScheduler(taskRepo repository.JobTaskRepo, logRepo repository.ExecLogRepo, rdb *redis.Client, logger *zap.Logger) *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		taskRepo: taskRepo,
		logRepo:  logRepo,
		rdb:      rdb,
		logger:   logger,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	tasks, err := s.taskRepo.ListActive(ctx, 0)
	if err != nil {
		return fmt.Errorf("加载活跃任务失败: %w", err)
	}

	for _, task := range tasks {
		if err := s.addTask(task); err != nil {
			s.logger.Error("注册定时任务失败", zap.String("task", task.Name), zap.Error(err))
		}
	}

	s.cron.Start()
	s.logger.Info("调度器已启动", zap.Int("taskCount", len(tasks)))
	return nil
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) addTask(task model.JobTask) error {
	jobFunc := s.buildJobFunc(task)
	_, err := s.cron.AddFunc(task.CronExpr, jobFunc)
	return err
}

func (s *Scheduler) buildJobFunc(task model.JobTask) func() {
	return func() {
		ctx := context.Background()
		now := time.Now().Format("20060102150405")
		lockKey := fmt.Sprintf("job:exec:%d:%s", task.ID, now)

		// 幂等保护：同时间只有一个实例执行
		locked, err := s.rdb.SetNX(ctx, lockKey, 1, 5*time.Minute).Result()
		if err != nil || !locked {
			return
		}

		logID := snowflake.NextID()
		execLog := &model.JobExecutionLog{
			ID:          logID,
			TenantID:    task.TenantID,
			TaskID:      task.ID,
			TaskName:    task.Name,
			TriggerType: model.TriggerTypeCron,
			Status:      model.ExecStatusRunning,
			StartTime:   time.Now(),
		}
		_ = s.logRepo.Create(ctx, execLog)

		start := time.Now()
		err = s.executeHandler(task.Handler, task.Params)
		duration := time.Since(start).Milliseconds()

		endTime := time.Now()
		updates := map[string]interface{}{
			"duration_ms": duration,
			"end_time":    endTime,
		}
		if err != nil {
			updates["status"] = model.ExecStatusFailed
			updates["error"] = err.Error()
		} else {
			updates["status"] = model.ExecStatusSuccess
		}
		_ = s.logRepo.Update(ctx, logID, updates)
	}
}

func (s *Scheduler) executeHandler(handler, params string) error {
	// 内置处理器注册表
	handlers := map[string]func(string) error{}
	if fn, ok := handlers[handler]; ok {
		return fn(params)
	}
	return fmt.Errorf("未注册的处理器: %s", handler)
}
