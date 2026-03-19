package portfolio

import (
	"quant-trading/internal/domain/order"
)

type Engine interface {
	UpdateFill(symbol string, side order.Side, price float64, qty int64)
	UpdateMarkPrice(symbol string, price float64)

	UnrealizedPnL() float64
	RealizedPnL() float64

	Snapshot() Snapshot
	Restore(Snapshot)
}
