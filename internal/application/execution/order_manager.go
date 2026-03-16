package execution

import (
	"fmt"
	"quant-trading/internal/domain/order"
	"sync"
)

type OrderManager struct {
	mu     sync.Mutex
	orders map[string]*order.Order
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders: make(map[string]*order.Order),
	}
}

func (om *OrderManager) CreateOrder(
	accountID, strategyID, symbol string,
	side order.Side, orderType order.OrderType, qty float64, price float64,
) *order.Order {
	om.mu.Lock()
	defer om.mu.Unlock()

	id := fmt.Sprintf("%s-%d", symbol, len(om.orders)+1)
	o := order.NewOrder(id, strategyID, accountID, symbol, side, orderType, price, qty)
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
