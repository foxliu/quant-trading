package strategy

import (
	"context"
	"quant-trading/internal/domain"
)

/*
最小定义
*/

type DemoStrategy struct {
	name string
}

func NewDemoStrategy(name string) *DemoStrategy {
	return &DemoStrategy{
		name: name,
	}
}

func (d *DemoStrategy) Name() string {
	return d.name
}

func (d *DemoStrategy) OnBar(ctx context.Context, sc Context, bar domain.MarketBar) *domain.Signal {
	if bar.Close > 100 {
		return &domain.Signal{
			Instrument: bar.Instrument,
			Side:       domain.Buy,
			Strength:   0.5,
			Time:       sc.Now(),
			Reason:     "demo close > 100",
		}
	}
	return nil
}
