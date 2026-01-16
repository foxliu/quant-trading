package trade

import "time"

/*
Position

========

Position 是“事实状态”，只由 Execution Event 推导

	Qty 永远是 正数

	Side 表示方向

	0 Qty = 无仓位（Position 不存在）
*/
type Position struct {
	AccountID string
	Symbol    string
	Sid       Side // Long / Short

	Qty int64 // 当前持仓数量（绝对值）

	AvgPrice float64 // 加权平均成本价

	OpenTime time.Time
	UpdateAt time.Time
}
