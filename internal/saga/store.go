package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// SagaInstance Saga 实例（MySQL 持久化）
type SagaInstance struct {
	ID          int64     `gorm:"primaryKey"`
	SagaName    string    `gorm:"column:saga_name;size:100;index"`
	State       string    `gorm:"column:state;size:20;index"`
	CurrentStep int       `gorm:"column:current_step"`
	Payload     string    `gorm:"column:payload;type:text"`
	StepResults string    `gorm:"column:step_results;type:text"`
	Error       string    `gorm:"column:error;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (SagaInstance) TableName() string { return "saga_instance" }

// SagaStore MySQL 状态存储
type SagaStore struct {
	db *gorm.DB
}

// NewSagaStore 创建存储
func NewSagaStore(db *gorm.DB) *SagaStore {
	return &SagaStore{db: db}
}

// Create 创建 Saga 实例
func (s *SagaStore) Create(ctx context.Context, instance *SagaInstance) error {
	return s.db.WithContext(ctx).Create(instance).Error
}

// Update 更新状态
func (s *SagaStore) Update(ctx context.Context, instance *SagaInstance) error {
	return s.db.WithContext(ctx).Save(instance).Error
}

// GetByID 按 ID 查询
func (s *SagaStore) GetByID(ctx context.Context, id int64) (*SagaInstance, error) {
	var inst SagaInstance
	if err := s.db.WithContext(ctx).First(&inst, id).Error; err != nil {
		return nil, err
	}
	return &inst, nil
}

// StepResult 步骤执行结果
type StepResult struct {
	StepName string `json:"stepName"`
	Output   string `json:"output"`
}

// SaveStepResult 保存步骤结果
func (s *SagaStore) SaveStepResult(ctx context.Context, inst *SagaInstance, result StepResult) error {
	var results []StepResult
	if inst.StepResults != "" {
		_ = json.Unmarshal([]byte(inst.StepResults), &results)
	}
	results = append(results, result)
	data, _ := json.Marshal(results)
	inst.StepResults = string(data)
	return s.Update(ctx, inst)
}

// AutoMigrate 自动迁移
func (s *SagaStore) AutoMigrate() error {
	return s.db.AutoMigrate(&SagaInstance{})
}

// DeadLetterSave 保存失败 Saga 到死信记录
func (s *SagaStore) DeadLetterSave(ctx context.Context, inst *SagaInstance, reason string) error {
	inst.State = string(StateFailed)
	inst.Error = fmt.Sprintf("DLQ: %s", reason)
	return s.Update(ctx, inst)
}
