// Package domain trade 成交
package trade

import (
	"quant-trading/internal/domain/common"
	"quant-trading/internal/domain/instrument"
	"time"
)

// Trade 成交
type Trade struct {
	OrderID    string
	Instrument instrument.Instrument

	Side     common.Side
	Price    float64
	Quantity float64
	Time     time.Time
}
