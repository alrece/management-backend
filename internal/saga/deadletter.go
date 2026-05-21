package saga

import (
	"context"

	"go.uber.org/zap"
)

// DeadLetterHandler 死信队列处理器
type DeadLetterHandler struct {
	store  *SagaStore
	logger *zap.Logger
}

// NewDeadLetterHandler 创建死信处理器
func NewDeadLetterHandler(store *SagaStore, logger *zap.Logger) *DeadLetterHandler {
	return &DeadLetterHandler{store: store, logger: logger}
}

// ListFailed 查询所有失败实例
func (h *DeadLetterHandler) ListFailed(ctx context.Context) ([]SagaInstance, error) {
	var list []SagaInstance
	err := h.store.db.WithContext(ctx).
		Where("state = ?", string(StateFailed)).
		Order("updated_at DESC").
		Find(&list).Error
	return list, err
}

// Retry 重试失败的 Saga
func (h *DeadLetterHandler) Retry(ctx context.Context, id int64, orchestrator *Orchestrator) error {
	inst, err := h.store.GetByID(ctx, id)
	if err != nil {
		return err
	}
	h.logger.Info("重试 Saga", zap.Int64("id", id), zap.String("saga", inst.SagaName))
	_, err = orchestrator.Execute(ctx, inst.SagaName, []byte(inst.Payload))
	return err
}
