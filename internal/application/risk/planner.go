package risk

// FROZEN: V1

import (
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/strategy"
)

// 从意图(intent)到订单(Orders)的推进

/*
Planner
=======

Planner 的工程职责（当前阶段）：

- 接收 Strategy Signal
- 将 Signal 翻译为标准化 Order
- 不做风控
- 不做仓位计算
- 不阻塞上游
*/
type Planner struct {
	output chan<- order.Order
}

func NewPlanner(output chan<- order.Order) *Planner {
	return &Planner{
		output: output,
	}
}

func (p *Planner) Consume(signal strategy.Signal) {
	orders := p.plan(signal)
	for _, o := range orders {
		select {
		case p.output <- o:
		default:
			// planner 不阻塞上游， buffer 满直接丢弃
		}
	}
}

func (p *Planner) plan(signal strategy.Signal) []order.Order {
	switch signal.Intent {
	case strategy.IntentLong:
		return p.openLong(signal)

	case strategy.IntentShort:
		return p.openShort(signal)

	case strategy.IntentFlat:
		return p.closePosition(signal)

	default:
		return nil
	}
}
