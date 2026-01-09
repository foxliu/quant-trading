package strategy

import (
	"context"
	"quant-trading/internal/domain"
)

/*
最小定义
*/

type Engine struct {
	strategies []Strategy
}

func NewEngine(strategies ...Strategy) *Engine {
	return &Engine{
		strategies: strategies,
	}
}

func (e *Engine) RunOnBar(ctx context.Context, sc Context, bar domain.MarketBar) []*domain.Signal {
	var signals []*domain.Signal

	for _, s := range e.strategies {
		sig := s.OnBar(ctx, sc, bar)
		if sig == nil {
			continue
		}
		signals = append(signals, sig)
	}
	return signals
}
