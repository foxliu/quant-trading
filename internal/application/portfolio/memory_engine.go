/*
Package portfolio

CAUTION: TEMPORARY PACKAGE
Deprecated: 仅用于Paper/Backtest
*/
package portfolio

import (
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/portfolio"
	"sync"
)

// MemoryEngine 内存实现 portfolio.Engine
type MemoryEngine struct {
	mu                sync.RWMutex
	positions         map[string]portfolio.Position
	peakEquity        float64
	dailyPnL          float64
	currentTradingDay string
}

func NewMemoryEngine() portfolio.Engine {
	return &MemoryEngine{
		positions:         make(map[string]portfolio.Position),
		peakEquity:        0,
		currentTradingDay: "",
	}
}

// UpdateTradingDay 交易日更新
func (e *MemoryEngine) UpdateTradingDay(tradingDay string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if tradingDay != e.currentTradingDay {
		e.currentTradingDay = tradingDay
		e.dailyPnL = 0
	}
}

func (e *MemoryEngine) UpdateFill(symbol string, side order.Side, price float64, qty int64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	p, exists := e.positions[symbol]
	if !exists {
		p = portfolio.Position{Instrument: instrument.Instrument{Symbol: symbol}}
	}
	if side == order.Buy {
		p.Quantity += qty
		if p.Quantity != 0 {
			p.OpenPrice = (p.OpenPrice*float64(p.Quantity-qty) + price/float64(qty)) / float64(p.Quantity)
		}
	} else {
		p.Quantity -= qty
	}
	e.positions[symbol] = p
}

func (e *MemoryEngine) UpdateMarkPrice(symbol string, price float64) {

}

func (e *MemoryEngine) UnrealizedPnL() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var total float64
	for _, p := range e.positions {
		total += p.UnrealizedPnL()
	}
	return total
}

func (e *MemoryEngine) RealizedPnL() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.dailyPnL
}

func (e *MemoryEngine) Snapshot() portfolio.Snapshot {
	return portfolio.Snapshot{}
}

func (e *MemoryEngine) Restore(s portfolio.Snapshot) {

}

func (e *MemoryEngine) GetPositions() ([]portfolio.Position, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	pos := make([]portfolio.Position, 0, len(e.positions))
	for _, p := range e.positions {
		pos = append(pos, p)
	}
	return pos, nil
}

func (e *MemoryEngine) GetPosition(symbol string) (portfolio.Position, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	p, ok := e.positions[symbol]
	return p, ok
}

func (e *MemoryEngine) GetDailyRealizedPnL() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.dailyPnL
}

func (e *MemoryEngine) GetMaxDrawdown() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	// TODO: 真实计算当前回撤（需要peakEquity和当前权益）
	return 0.0
}
