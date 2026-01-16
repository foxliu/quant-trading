package order

import (
	"quant-trading/internal/domain/strategy"
	"quant-trading/internal/domain/trade"
	"time"
)

/*
Order
=====

Order 表示一笔【可执行的标准化交易意图】。

约定：
- Order 是 Risk Engine 的最终产物
- Execution Engine 必须能够直接消费 Order
- Order 不表达策略意图（IntentLong / IntentShort）
- 只表达交易动作（Buy / Sell）
- Quantity 是相对变化量（Delta）
- Intent 复用 strategy.PositionIntent
- Side 使用 trade.Side
*/
type Order struct {
	// === 标识 ===
	OrderID    string
	StrategyID string
	Symbol     string

	// === 方向与行为 ===
	Side   trade.Side              // Buy / Sell
	Intent strategy.PositionIntent // IntentLong / IntentShort / IntentFlat

	// === 数量与价格 ===
	Quantity int64   // 相对数量 (Delta)
	Price    float64 // 0 = 市价

	// === 状态 ===
	Status Status

	// === 时间 ===
	CreatedAt time.Time
}
