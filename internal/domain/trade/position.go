package trade

import (
	"quant-trading/internal/domain/instrument"
	"time"
)

/*
Position

========
Qty 含义：

- Qty > 0 : 净多头

- Qty < 0 : 净空头

- Qty = 0 : 无仓位（Position 不存在）
*/
type Position struct {
	Instrument instrument.Instrument
	Symbol     string

	Qty int64 // 当前持仓数量（绝对值）

	AvgPrice float64 // 加权平均成本价

	LastPrice float64

	OpenTime time.Time
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

func (p *Position) UnrealizedPnL() float64 {
	return float64(p.Qty) * (p.LastPrice - p.AvgPrice)
}
