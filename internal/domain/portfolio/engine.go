package portfolio

import "quant-trading/internal/domain/trade"

type Engine interface {
	UpdateFill(symbol string, side trade.Side, qty int64, price float64)
	UpdateMarkPrice(symbol string, price float64)

	UnrealizedPnL() float64
	RealizedPnL() float64

	Snapshot() Snapshot
	Restore(Snapshot)
}
