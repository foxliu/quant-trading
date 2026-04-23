// Package capital
//
// CAUTION: TEMPORARY PACKAGE
// Deprecated: 仅用于Paper/Backtest
package capital

import (
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/order"
	"sync"
)

// MemoryEngine 内存实现 capital.Engine（paper/backtest 共用）
type MemoryEngine struct {
	mu        sync.RWMutex
	available float64
	frozen    float64
	total     float64
}

func NewMemoryEngine(initial float64) capital.Engine {
	return &MemoryEngine{
		available: initial,
		frozen:    0,
		total:     initial,
	}
}

func (e *MemoryEngine) Freeze(orderID, symbol string, price, qty float64, side order.Side) error {
	amount := price * qty
	if e.available < amount {
		return nil // paper宽松模式
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.available -= amount
	e.frozen += amount
	return nil
}

func (e *MemoryEngine) Commit(orderID string, amount float64) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.frozen -= amount
	return nil
}

func (e *MemoryEngine) Release(orderID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.frozen = 0
	return nil
}

func (e *MemoryEngine) Available() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.available
}
func (e *MemoryEngine) Total() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.total
}
func (e *MemoryEngine) Frozen() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.frozen
}
func (e *MemoryEngine) Snapshot() capital.Snapshot {
	return capital.Snapshot{}
}
func (e *MemoryEngine) Restore(s capital.Snapshot) {
}
