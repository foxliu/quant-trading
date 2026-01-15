package trade

import "time"

/*
Position
========

Position 表示某一 Symbol 在某一账户下的净持仓。

约定：
- Qty > 0  : 多头
- Qty < 0  : 空头
- Qty == 0 : 空仓
*/
type Position struct {
	AccountID string
	Symbol    string

	Qty float64 // 当前净仓位(+Long/-Short)

	AvgPrice float64 // 持仓均价 (当前阶段不参与计算)

	UpdateAt time.Time
}

func (p *Position) IsFlat() bool {
	return p.Qty == 0
}

func (p *Position) IsLong() bool {
	return p.Qty > 0
}

func (p *Position) IsShort() bool {
	return p.Qty < 0
}

func (p *Position) AbsQty() float64 {
	if p.Qty < 0 {
		return -p.Qty
	}
	return p.Qty
}
