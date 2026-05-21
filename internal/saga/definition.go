package saga

import "context"

// SagaState 状态机
type SagaState string

const (
	StatePending      SagaState = "PENDING"
	StateRunning      SagaState = "RUNNING"
	StateCompensating SagaState = "COMPENSATING"
	StateCompleted    SagaState = "COMPLETED"
	StateFailed       SagaState = "FAILED"
)

// StepDefinition Saga 步骤定义
type StepDefinition struct {
	Name       string
	Execute    func(ctx context.Context, payload []byte) ([]byte, error)
	Compensate func(ctx context.Context, payload []byte) error
}

// Definition Saga 定义 DSL
type Definition struct {
	Name  string
	Steps []StepDefinition
}

// NewDefinition 创建 Saga 定义
func NewDefinition(name string) *Definition {
	return &Definition{Name: name}
}

// AddStep 添加步骤
func (d *Definition) AddStep(name string, execute func(ctx context.Context, payload []byte) ([]byte, error), compensate func(ctx context.Context, payload []byte) error) *Definition {
	d.Steps = append(d.Steps, StepDefinition{
		Name:       name,
		Execute:    execute,
		Compensate: compensate,
	})
	return d
}
