package risk

import (
	"quant-trading/internal/domain/order"
)

type Decision struct {
	Orders []order.Order
}
