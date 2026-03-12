package risk

// FROZEN: V1

import (
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/strategy"
	"quant-trading/internal/domain/trade"
	"time"
)

/*
Execution Planning Helpers
==========================

约定：
- Signal.TargetQty 表示“期望最终仓位规模（绝对值）”
- Position.Qty 表示“当前净仓位（可正可负）”
- Order.Qty 永远是正数
*/

// 开多仓
func (p *Planner) openLong(signal strategy.Signal) []execution.Order {
	o := execution.Order{
		StrategyID: signal.StrategyID,
		Symbol:     signal.Symbol,
		Side:       trade.Buy,
		Intent:     strategy.IntentLong,
		Quantity:   int64(signal.TargetQty), // 此处 Quantity 暂时等于 TargetQty 后续会由 Position Engine 修正为 Delta
		Price:      signal.Price,
		Status:     execution.Pending,
		CreatedAt:  time.Now(),
	}
	return []execution.Order{o}
}

// 开空仓
func (p *Planner) openShort(signal strategy.Signal) []execution.Order {
	o := execution.Order{
		StrategyID: signal.StrategyID,
		Symbol:     signal.Symbol,

		Side:   trade.Sell,
		Intent: strategy.IntentShort,

		Quantity: int64(signal.TargetQty), //此处 Quantity 暂时等于 TargetQty 后续会由 Position Engine 修正为 Delta
		Price:    signal.Price,

		Status:    execution.Pending,
		CreatedAt: time.Now(),
	}
	return []execution.Order{o}
}

// 平仓
func (p *Planner) closePosition(signal strategy.Signal) []execution.Order {
	o := execution.Order{
		OrderID:    "",
		StrategyID: signal.StrategyID,
		Symbol:     signal.Symbol,
		Side:       trade.Sell, // 占位，后续由 Position 决定
		Intent:     strategy.IntentFlat,
		Quantity:   int64(signal.TargetQty), // 此处 Quantity 暂时等于 TargetQty 后续会由 Position Engine 修正为 Delta
		Price:      signal.Price,
		Status:     execution.Pending,
		CreatedAt:  time.Now(),
	}
	return []execution.Order{o}
}
