package portfolio

import (
	"quant-trading/internal/domain/order"
	"quant-trading/pkg/utils"
	"sync"
)

type position struct {
	symbol    string
	quantity  float64
	avgPrice  float64
	lastPrice float64
}

type SimplePortfolio struct {
	mu sync.RWMutex

	positions map[string]*position
	realized  float64
}

func NewSimplePortfolio() *SimplePortfolio {
	return &SimplePortfolio{
		positions: make(map[string]*position),
	}
}

func (p *SimplePortfolio) UpdateFill(symbol string, side order.Side, price float64, qty float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pos, ok := p.positions[symbol]
	if !ok {
		pos = &position{symbol: symbol}
		p.positions[symbol] = pos
	}

	signedQty := qty
	if side == order.Sell {
		signedQty = -qty
	}

	newQty := pos.quantity + signedQty

	if pos.quantity == 0 {
		pos.avgPrice = price
	} else if (pos.quantity > 0 && signedQty > 0) || (pos.quantity < 0 && signedQty < 0) {
		totalCost := pos.quantity*pos.avgPrice + signedQty*price
		pos.avgPrice = totalCost / newQty
	} else {
		closeQty := min(utils.Abs(pos.quantity), utils.Abs(signedQty))
		p.realized += closeQty * (price - pos.avgPrice) * utils.Sign(pos.quantity)
	}
	pos.quantity = newQty
}

func (p *SimplePortfolio) UpdateMarkPrice(symbol string, price float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if pos, ok := p.positions[symbol]; ok {
		pos.lastPrice = price
	}
}

func (p *SimplePortfolio) UnrealizedPnL() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := 0.0
	for _, pos := range p.positions {
		total += pos.quantity * (pos.lastPrice - pos.avgPrice)
	}
	return total
}

func (p *SimplePortfolio) RealizedPnL() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.realized
}

func (p *SimplePortfolio) Snapshot() Snapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var snaps []PositionSnapshot
	for _, pos := range p.positions {
		snaps = append(snaps, PositionSnapshot{
			Symbol:    pos.symbol,
			Quantity:  pos.quantity,
			AvgPrice:  pos.avgPrice,
			LastPrice: pos.lastPrice,
		})
	}
	return Snapshot{
		Positions: snaps,
		Realized:  p.realized,
	}
}

func (p *SimplePortfolio) Restore(s Snapshot) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.positions = make(map[string]*position)
	for _, ps := range s.Positions {
		p.positions[ps.Symbol] = &position{
			symbol:    ps.Symbol,
			quantity:  ps.Quantity,
			avgPrice:  ps.AvgPrice,
			lastPrice: ps.LastPrice,
		}
	}
	p.realized = s.Realized
}
