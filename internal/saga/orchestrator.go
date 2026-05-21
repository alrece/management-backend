package saga

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Orchestrator Saga 编排器（Redis Streams 事件驱动）
type Orchestrator struct {
	store  *SagaStore
	rdb    *redis.Client
	logger *zap.Logger
	defs   map[string]*Definition
}

// NewOrchestrator 创建编排器
func NewOrchestrator(store *SagaStore, rdb *redis.Client, logger *zap.Logger) *Orchestrator {
	return &Orchestrator{
		store:  store,
		rdb:    rdb,
		logger: logger,
		defs:   make(map[string]*Definition),
	}
}

// Register 注册 Saga 定义
func (o *Orchestrator) Register(def *Definition) {
	o.defs[def.Name] = def
}

// Execute 执行 Saga
func (o *Orchestrator) Execute(ctx context.Context, sagaName string, payload []byte) (int64, error) {
	def, ok := o.defs[sagaName]
	if !ok {
		return 0, fmt.Errorf("Saga 定义不存在: %s", sagaName)
	}

	inst := &SagaInstance{
		SagaName: sagaName,
		State:    string(StateRunning),
		Payload:  string(payload),
	}

	if err := o.store.Create(ctx, inst); err != nil {
		return 0, err
	}

	go o.run(context.Background(), def, inst)

	return inst.ID, nil
}

func (o *Orchestrator) run(ctx context.Context, def *Definition, inst *SagaInstance) {
	var executedSteps []int

	for i, step := range def.Steps {
		inst.CurrentStep = i
		inst.State = string(StateRunning)
		_ = o.store.Update(ctx, inst)

		output, err := step.Execute(ctx, []byte(inst.Payload))
		if err != nil {
			o.logger.Error("Saga 步骤执行失败",
				zap.String("saga", def.Name),
				zap.String("step", step.Name),
				zap.Error(err),
			)

			o.compensate(ctx, def, inst, executedSteps, err)
			return
		}

		executedSteps = append(executedSteps, i)
		_ = o.store.SaveStepResult(ctx, inst, StepResult{
			StepName: step.Name,
			Output:   string(output),
		})

		// 发布事件到 Redis Stream
		o.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: fmt.Sprintf("saga:%s:events", def.Name),
			Values: map[string]interface{}{
				"instance_id": inst.ID,
				"step":        step.Name,
				"state":       "completed",
			},
		})
	}

	inst.State = string(StateCompleted)
	_ = o.store.Update(ctx, inst)
	o.logger.Info("Saga 完成", zap.String("saga", def.Name), zap.Int64("id", inst.ID))
}

func (o *Orchestrator) compensate(ctx context.Context, def *Definition, inst *SagaInstance, executedSteps []int, execErr error) {
	inst.State = string(StateCompensating)
	inst.Error = execErr.Error()
	_ = o.store.Update(ctx, inst)

	// 逆序执行补偿
	for i := len(executedSteps) - 1; i >= 0; i-- {
		step := def.Steps[executedSteps[i]]
		if step.Compensate != nil {
			if err := step.Compensate(ctx, []byte(inst.Payload)); err != nil {
				o.logger.Error("Saga 补偿失败",
					zap.String("saga", def.Name),
					zap.String("step", step.Name),
					zap.Error(err),
				)
				_ = o.store.DeadLetterSave(ctx, inst, fmt.Sprintf("补偿失败: step=%s err=%v", step.Name, err))
				return
			}
		}
	}

	_ = o.store.DeadLetterSave(ctx, inst, execErr.Error())
}
