package risk

import (
	"quant-trading/internal/domain/execution"
)

type Decision struct {
	Orders []execution.Order
}
