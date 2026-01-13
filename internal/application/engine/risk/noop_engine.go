package risk

import (
	"context"
	"quant-trading/internal/domain/strategy"
)

/*
NoopEngine
==========

用于工程打通阶段的占位实现：
- 吃掉所有 Signal
- 不做任何处理
*/
type NoopEngine struct{}

func NewNoopEngine() *NoopEngine {
	return &NoopEngine{}
}

func (n *NoopEngine) Start(ctx context.Context) error {
	return nil
}

func (n *NoopEngine) Stop() error {
	return nil
}

func (n *NoopEngine) Consume(signal strategy.Signal) {
	// 故意留空
}
