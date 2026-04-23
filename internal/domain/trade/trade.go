// Package trade  成交
package trade

import (
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/order"
	"time"
)

// Trade 成交
type Trade struct {
	OrderID    string
	Instrument instrument.Instrument

	Side     order.Side
	Price    float64
	Quantity float64
	Time     time.Time
}
