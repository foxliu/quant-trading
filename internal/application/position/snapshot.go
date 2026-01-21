package position

import "quant-trading/internal/domain/trade"

// 只读出口

type Snapshot struct {
	Symbol   string
	Side     trade.Side
	Qty      int64
	AvgPrice float64
}

func NewSnapshot(pos *trade.Position) *Snapshot {
	if pos == nil {
		return nil
	}

	return &Snapshot{
		Symbol:   pos.Symbol,
		Qty:      pos.Qty,
		AvgPrice: pos.AvgPrice,
	}
}
