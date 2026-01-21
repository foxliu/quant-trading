package pnl

import "quant-trading/internal/domain/trade"

/*
Unrealized PnL
==============
Qty 正负自然表达方向
*/
func Unrealized(pos *trade.Position, marketPrice float64) float64 {
	if pos == nil || pos.Qty == 0 {
		return 0
	}

	return float64(pos.Qty) * (marketPrice - pos.AvgPrice)
}
