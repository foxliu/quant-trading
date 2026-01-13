package risk

import (
	"context"
	"quant-trading/internal/domain/strategy"
)

/*
Engine
======

Risk Engine 的工程职责（当前阶段）：
- 接收 Strategy 产生的 Signal
- 不做任何决策
- 不阻塞上游
*/
type Engine interface {
	Start(ctx context.Context) error
	Stop() error

	Consume(signal strategy.Signal)
}
