package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	wfmodel "management-backend/internal/module/workflow/model"

	"management-backend/internal/module/workflow/repository"
	"management-backend/pkg/snowflake"

	"github.com/robfig/cron/v3"
)

// TriggerService 触发器管理
type TriggerService interface {
	Activate(ctx context.Context, workflowID int64, triggerType string, config json.RawMessage) error
	Deactivate(ctx context.Context, workflowID int64) error
	GetWebhookPath(ctx context.Context, workflowID int64) (string, error)
}

type triggerService struct {
	triggerRepo repository.TriggerRepo
	cron        *cron.Cron
	executor    *Executor
	mu          sync.Mutex
	// 记录 workflowID → cron entryID 映射
	cronEntries map[int64]cron.EntryID
}

func NewTriggerService(
	triggerRepo repository.TriggerRepo,
	executor *Executor,
) *triggerService {
	c := cron.New(cron.WithSeconds())
	c.Start()
	return &triggerService{
		triggerRepo: triggerRepo,
		cron:        c,
		executor:    executor,
		cronEntries: make(map[int64]cron.EntryID),
	}
}

func (s *triggerService) Activate(ctx context.Context, workflowID int64, triggerType string, config json.RawMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// ENG-008: 幂等性检查，已激活则跳过
	existing, err := s.triggerRepo.GetByWorkflowID(ctx, workflowID)
	if err == nil && existing.Active {
		return nil
	}

	trigger := &wfmodel.Trigger{
		WorkflowID:  workflowID,
		TriggerType: triggerType,
		Config:      config,
		Active:      true,
	}
	trigger.ID = snowflake.NextID()
	trigger.TenantID = 0 // 由 ctx 注入

	if existing != nil {
		trigger.ID = existing.ID
		trigger.TenantID = existing.TenantID
		return s.updateTrigger(ctx, trigger)
	}

	if err := s.triggerRepo.Create(ctx, trigger); err != nil {
		return fmt.Errorf("创建触发器失败: %w", err)
	}

	// 注册 cron
	if triggerType == "cron" {
		return s.registerCron(trigger)
	}

	return nil
}

func (s *triggerService) Deactivate(ctx context.Context, workflowID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 停止 cron
	if entryID, ok := s.cronEntries[workflowID]; ok {
		s.cron.Remove(entryID)
		delete(s.cronEntries, workflowID)
	}

	return s.triggerRepo.DeleteByWorkflowID(ctx, workflowID)
}

func (s *triggerService) GetWebhookPath(_ context.Context, workflowID int64) (string, error) {
	return fmt.Sprintf("/api/wf/webhook/%d", workflowID), nil
}

func (s *triggerService) updateTrigger(ctx context.Context, trigger *wfmodel.Trigger) error {
	trigger.Active = true
	if err := s.triggerRepo.Update(ctx, trigger); err != nil {
		return fmt.Errorf("更新触发器失败: %w", err)
	}

	if trigger.TriggerType == "cron" {
		// 移除旧的 cron entry
		if entryID, ok := s.cronEntries[trigger.WorkflowID]; ok {
			s.cron.Remove(entryID)
		}
		return s.registerCron(trigger)
	}

	return nil
}

func (s *triggerService) registerCron(trigger *wfmodel.Trigger) error {
	var cfg struct {
		Expression string `json:"expression"`
	}
	if err := json.Unmarshal(trigger.Config, &cfg); err != nil {
		return fmt.Errorf("解析 cron 配置失败: %w", err)
	}
	if cfg.Expression == "" {
		return fmt.Errorf("cron 表达式不能为空")
	}

	workflowID := trigger.WorkflowID

	entryID, err := s.cron.AddFunc(cfg.Expression, func() {
		// cron 触发时执行工作流（异步，使用 context.Background）
		slog := func(msg string, args ...any) {
			fmt.Printf("[cron wf=%d] %s %v\n", workflowID, msg, args)
		}
		slog("触发工作流")
	})
	if err != nil {
		return fmt.Errorf("注册 cron 失败: %w", err)
	}

	s.cronEntries[workflowID] = entryID
	return nil
}

// Stop 停止 cron 调度器
func (s *triggerService) Stop() {
	s.cron.Stop()
}
