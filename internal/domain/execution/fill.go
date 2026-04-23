package execution

import (
	"quant-trading/internal/domain/order"
	"time"
)

type Fill struct {
	OrderID   string
	AccountID string

	Symbol string
	Side   order.Side

	Qty   int64
	Price float64
	Fee   float64

	Timestamp time.Time
}
