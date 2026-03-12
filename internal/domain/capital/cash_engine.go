package capital

import (
	"errors"
	"quant-trading/internal/domain/execution"
	"sync"
)

var ErrInsufficientCapital = errors.New("insufficient capital")

type CashEngine struct {
	mu sync.Mutex

	total     float64
	available float64
	frozen    float64

	frozenOrders map[string]float64
}

func NewCashEngine(initial float64) *CashEngine {
	return &CashEngine{
		total:        initial,
		available:    initial,
		frozen:       0,
		frozenOrders: make(map[string]float64),
	}
}

func (e *CashEngine) Freeze(orderID, symbol string, price float64, qty float64, side execution.Side) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	required := price * qty

	if e.available < required {
		return ErrInsufficientCapital
	}

	e.available -= required
	e.frozen += required
	e.frozenOrders[orderID] = required
	return nil
}

func (e *CashEngine) Commit(orderID string, amount float64) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	frozen, ok := e.frozenOrders[orderID]
	if !ok {
		return nil
	}

	if amount > frozen {
		amount = frozen
	}

	e.frozen -= amount
	e.total -= amount

	remaining := frozen - amount
	if remaining == 0 {
		delete(e.frozenOrders, orderID)
	} else {
		e.frozenOrders[orderID] = remaining
	}
	return nil
}

func (e *CashEngine) Release(orderID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	frozen, ok := e.frozenOrders[orderID]
	if !ok {
		return nil
	}

	e.frozen -= frozen
	e.available += frozen

	delete(e.frozenOrders, orderID)
	return nil
}

func (e *CashEngine) Available() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.available
}

func (e *CashEngine) Frozen() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.frozen
}

func (e *CashEngine) Total() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.total
}

func (e *CashEngine) Snapshot() Snapshot {
	e.mu.Lock()
	defer e.mu.Unlock()

	copyMap := make(map[string]float64)
	for k, v := range e.frozenOrders {
		copyMap[k] = v
	}

	return Snapshot{
		Total:        e.total,
		Available:    e.available,
		Frozen:       e.frozen,
		FrozenOrders: copyMap,
	}
}

func (e *CashEngine) Restore(s Snapshot) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.total = s.Total
	e.available = s.Available
	e.frozen = s.Frozen
	e.frozenOrders = make(map[string]float64)
	for k, v := range s.FrozenOrders {
		e.frozenOrders[k] = v
	}
}
