package trade

import "time"

// Position 是执行链路中使用的最小仓位事实模型。
// Qty 正负表示方向：正=多头，负=空头。
type Position struct {
	Symbol   string
	Qty      int64
	AvgPrice float64
	OpenTime time.Time
	UpdateAt time.Time
}
