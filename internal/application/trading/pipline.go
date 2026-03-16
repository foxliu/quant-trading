package trading

import (
	"context"
	"quant-trading/internal/application/position"
	"quant-trading/internal/application/risk"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/strategy"
)

/*
Pipeline
========

Pipeline 负责将交易系统的各个处理阶段串联起来。
*/
type Pipeline struct {
	// 输入
	signalCh <-chan strategy.Signal

	// 内部通道
	orderCh1 chan order.Order // planner -> Position
	orderCh2 chan order.Order // Position -> Risk
	orderCh3 chan order.Order // Risk -> Execution (保留）

	// 组件
	planner *risk.Planner
	// positionEngine *position.Engine
	riskEngine risk.Engine
}

func NewPipeline(
	signalCh <-chan strategy.Signal,
	positionCtx *position.Context,
	riskCtx *risk.Context,
) *Pipeline {
	orderCh1 := make(chan order.Order, 1024)
	orderCh2 := make(chan order.Order, 1024)
	orderCh3 := make(chan order.Order, 1024)

	planner := risk.NewPlanner(orderCh1)

	// TODO: 实现position engine和risk engine
	// var posEngine *position.Engine = nil
	var riskEngine risk.Engine = nil
	// posEngine := position.NewEngine(positionCtx, orderCh1, orderCh2)
	// riskEngine := risk.NewEngine(riskCtx, orderCh2, orderCh3)
	_ = positionCtx
	_ = riskCtx
	_ = orderCh2
	_ = orderCh3

	return &Pipeline{
		signalCh: signalCh,
		orderCh1: orderCh1,
		orderCh2: orderCh2,
		orderCh3: orderCh3,
		planner:  planner,
		// positionEngine: posEngine,
		riskEngine: riskEngine,
	}
}

func (p *Pipeline) Start(ctx context.Context) error {
	// 启动 Position / Risk (它们自己监听 channel）
	// _ = p.positionEngine.Start(ctx)
	if p.riskEngine != nil {
		_ = p.riskEngine.Start(ctx)
	}

	// 启动 Planner 消费 Signal
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-p.signalCh:
				p.planner.Consume(s)
			}
		}
	}()
	return nil
}
