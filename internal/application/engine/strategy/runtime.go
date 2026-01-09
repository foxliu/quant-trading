package strategyengine

import (
	"quant-trading/internal/domain/market"
	"quant-trading/internal/domain/strategy"
)

/*
Runtime
=======

Runtime 表示【一个策略实例的运行时】。

设计原则：
- 一个 Runtime = 一个 Strategy
- Runtime 内部串行执行
- 并发隔离在 Dispatcher 层完成
*/
type Runtime struct {
	strategy strategy.Strategy
	ctx      strategy.Context
}

func NewRuntime(s strategy.Strategy) *Runtime {
	return &Runtime{
		strategy: s,
		ctx:      strategy.NewContext(),
	}
}

func (r *Runtime) Init() error {
	return r.strategy.OnInit(r.ctx)
}

func (r *Runtime) Stop() error {
	return r.strategy.OnStop(r.ctx)
}

func (r *Runtime) HandleEvent(event market.Event) ([]strategy.Signal, error) {
	return r.strategy.OnMarketEvent(r.ctx, event)
}
