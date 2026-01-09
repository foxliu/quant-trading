package strategy

import (
	"context"
	"quant-trading/internal/domain"
)

/*
最小定义
*/

type Runtime struct {
	strategy Strategy
	ctx      Context
}

func NewRuntime(strategy Strategy, ctx Context) *Runtime {
	return &Runtime{
		strategy: strategy,
		ctx:      ctx,
	}
}

func (r *Runtime) OnBar(ctx context.Context, bar domain.MarketBar) *domain.Signal {
	return r.strategy.OnBar(ctx, r.ctx, bar)
}
