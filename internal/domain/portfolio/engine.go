package portfolio

import (
	"quant-trading/internal/domain/execution"
)

type Engine interface {
	UpdateFill(symbol string, side execution.Side, price float64, qty int64)
	UpdateMarkPrice(symbol string, price float64)

	UnrealizedPnL() float64
	RealizedPnL() float64

	Snapshot() Snapshot
	Restore(Snapshot)
}
