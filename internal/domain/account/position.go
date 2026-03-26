package account

type Side int

const (
	Buy Side = iota
	Sell
)

type Position struct {
	symbol      string
	qty         int64
	avg         float64
	realizedPnL float64
}

func newPosition(symbol string) *Position {
	return &Position{
		symbol: symbol,
	}
}

func (p *Position) applyFill(fill Fill) error {
	if fill.Side == Buy {
		totalCost := p.avg*float64(p.qty) + fill.Price*float64(fill.Qty)
		p.qty += fill.Qty
		p.avg = totalCost / float64(p.qty)
		return nil
	}

	// Sell
	if p.qty < fill.Qty {
		return ErrInsufficientPosition
	}

	// 计算已实现盈亏
	pnl := (fill.Price - p.avg) * float64(fill.Qty)
	p.realizedPnL += pnl

	p.qty -= fill.Qty
	if p.qty == 0 {
		p.avg = 0
	}
	return nil
}

func (p *Position) rollbackFill(fill Fill) {
	if fill.Side == Buy {
		p.qty -= fill.Qty
		if p.qty <= 0 {
			p.qty = 0
			p.avg = 0
		}
		return
	}
	p.qty += fill.Qty
}

func (p *Position) view(marketPrice float64) PositionView {
	unrealized := (marketPrice - p.avg) * float64(p.qty)
	return PositionView{
		Symbol:        p.symbol,
		Qty:           p.qty,
		Avg:           p.avg,
		RealizedPnL:   p.realizedPnL,
		UnrealizedPnL: unrealized,
		MarketValue:   marketPrice * float64(p.qty),
	}
}

type PositionView struct {
	Symbol        string
	Qty           int64
	Avg           float64
	RealizedPnL   float64
	UnrealizedPnL float64
	MarketValue   float64
}
