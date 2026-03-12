package execution

import (
	"quant-trading/internal/domain/trade"
)

type OrderState string

const (
	StateNew             OrderState = "NEW"
	StateSubmitted       OrderState = "SUBMITTED"
	StateAccepted        OrderState = "ACCEPT"
	StatePartiallyFilled OrderState = "PARTIALLY_FILLED"
	StateFilled          OrderState = "FILLED"
	StateCanceled        OrderState = "CANCELED"
	StateRejected        OrderState = "REJECTED"
)

type Order struct {
	ID     string
	Symbol string

	Side trade.Side

	Qty int64

	FilledQty int64

	State OrderState
}
