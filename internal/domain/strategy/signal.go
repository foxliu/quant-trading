package strategy

import (
	"time"
)

/*
Signal
======

Signal 表示【策略层发出的交易意图】。

工程原则：
- Signal 是“意图”，不是“指令”
- 不保证一定被执行
- 必须可被风控 / 仓位 / 下单模块独立处理
*/
type Signal struct {
	// === 基本标识 ===

	StrategyID string    // 产生该 Signal 的策略标识
	Symbol     string    // 交易标的（如 BTC-USDT / AAPL）
	Timestamp  time.Time // Signal 产生的时间

	// === 意图描述 ===

	Action Action         // Buy / Sell / Close / Hold
	Intent PositionIntent // Long / Short （对衍生品友好）

	// === 数量与价格（均为“建议”） ===

	Quantity float64 // 期望数量（非最终成交数量）
	Price    float64 // 期望价格（0 表示市价）

	// === 扩展字段 ===

	Meta map[string]string //策略自定义信息（不参与撮合）
}
