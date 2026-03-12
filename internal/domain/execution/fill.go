package execution

import "time"

type Fill struct {
	OrderID   string
	AccountID string

	Symbol string
	Side   Side

	Qty   int64
	Price float64
	Fee   float64

	Timestamp time.Time
}
