package position

//核心：Target → Delta

import (
	"quant-trading/internal/domain/strategy"
	"quant-trading/internal/domain/trade"
)

/*
CalcDelta
=========

根据当前仓位和策略 Signal，计算：

- 需要下多少量（Delta Qty）
- 最终交易方向（Buy / Sell）
*/
func CalcDelta(pos *trade.Position, signal strategy.Signal) (delta float64, side trade.Side) {
	current := 0.0
	if pos != nil {
		current = pos.Qty
	}

	target := signal.TargetQty
	delta = target - current

	switch {
	case delta > 0:
		side = trade.Buy
	case delta < 0:
		side = trade.Sell
	default:
		side = "" // 不需要交易
	}
	return
}
