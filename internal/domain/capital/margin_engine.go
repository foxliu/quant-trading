package capital

import (
	"quant-trading/internal/domain/trade"
)

// 暂时用不到，未来实现

type MarginEngine struct {
	leverage   float64
	marginRate float64
}

func (e *MarginEngine) Freeze(orderID, symbol string, price, qty float64, side trade.Side) error {
	// 未来实现
	return nil
}
