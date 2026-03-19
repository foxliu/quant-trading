/*
Package portfolio

CAUTION: TEMPORARY PACKAGE
Deprecated: 仅用于Paper/Backtest
*/
package portfolio

import (
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/portfolio"
	"sync"
)

// MemoryEngine 内存实现 portfolio.Engine
type MemoryEngine struct {
	mu sync.RWMutex

	positions map[string]float64
	avgPrices map[string]float64
}

func NewMemoryEngine() portfolio.Engine {
	return &MemoryEngine{
		positions: make(map[string]float64),
		avgPrices: make(map[string]float64),
	}
}

func (e *MemoryEngine) UpdateFill(symbol string, side order.Side, price float64, qty int64) {
	// 简化实现（与ApplyFill配合）
	qtyF := float64(qty)
	current := e.positions[symbol]
	if side == order.Buy {
		e.positions[symbol] = current + qtyF
		e.avgPrices[symbol] = (current*e.avgPrices[symbol] + qtyF*price) / (current + qtyF)
	} else {
		e.positions[symbol] = current - qtyF
	}
}

func (e *MemoryEngine) UpdateMarkPrice(symbol string, price float64) {

}

func (e *MemoryEngine) UnrealizedPnL() float64 {
	return 0
}

func (e *MemoryEngine) RealizedPnL() float64 {
	return 0
}

func (e *MemoryEngine) Snapshot() portfolio.Snapshot {
	return portfolio.Snapshot{}
}

func (e *MemoryEngine) Restore(s portfolio.Snapshot) {

}
