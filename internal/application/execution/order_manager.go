package execution

import (
	"fmt"
	"quant-trading/internal/domain/execution"
	"sync"
)

type OrderManager struct {
	mu     sync.Mutex
	orders map[string]*execution.Order
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders: make(map[string]*execution.Order),
	}
}

func (om *OrderManager) CreateOrder(
	accountID, strategyID, symbol string,
	side execution.Side, qty int64, price float64,
) *execution.Order {
	om.mu.Lock()
	defer om.mu.Unlock()

	id := fmt.Sprintf("%s-%d", symbol, len(om.orders)+1)
	o := execution.NewOrder(id, strategyID, accountID, symbol, side, price, qty)
	om.orders[id] = o

	return o
}

func (om *OrderManager) MarkFilled(orderID string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if o, ok := om.orders[orderID]; ok {
		o.MarkFilled()
	}
}
