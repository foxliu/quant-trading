package risk

import "quant-trading/internal/domain/trade"

type Context interface {
	Position(symbol string) trade.Position
	Equity() float64
}
