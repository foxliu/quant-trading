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
- 串行或并行消费（由实现决定）
- 内部进行风控 / 仓位 / 下单规划
- 不阻塞上游（Consume 必须是非阻塞的）
*/
type Engine interface {
	Start(ctx context.Context) error
	Stop() error

	Consume(signal strategy.Signal)
}
