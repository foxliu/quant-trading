package account

type Portfolio struct {
	positions map[string]*Position
}

func newPortfolio() *Portfolio {
	return &Portfolio{positions: make(map[string]*Position)}
}

func (p *Portfolio) applyFill(fill Fill) error {
	pos, ok := p.positions[fill.Symbol]
	if !ok {
		pos = newPosition(fill.Symbol)
		p.positions[fill.Symbol] = pos
	}
	return pos.applyFill(fill)
}

func (p *Portfolio) rollbackFill(fill Fill) {
	pos, ok := p.positions[fill.Symbol]
	if !ok {
		return
	}
	pos.rollbackFill(fill)
}

func (p *Portfolio) snapshot(market map[string]float64) (map[string]PositionView, float64, float64) {
	result := make(map[string]PositionView)

	var totalMarketValue float64
	var totalUnrealized float64

	for symbol, pos := range p.positions {
		price := market[symbol]

		view := pos.view(price)

		result[symbol] = view

		totalMarketValue += view.MarketValue
		totalUnrealized += view.UnrealizedPnL
	}
	return result, totalMarketValue, totalUnrealized
}
