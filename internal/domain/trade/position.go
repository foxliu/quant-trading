package trade

import "time"

/*
Position

========
Qty 含义：

- Qty > 0 : 净多头

- Qty < 0 : 净空头

- Qty = 0 : 无仓位（Position 不存在）
*/
type Position struct {
	Symbol string

	Qty int64 // 当前持仓数量（绝对值）

	AvgPrice float64 // 加权平均成本价

	OpenTime time.Time
	UpdateAt time.Time
}
