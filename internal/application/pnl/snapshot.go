package pnl

import "quant-trading/internal/domain/pnl"

func (c *Context) Snapshot() pnl.Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	return pnl.Snapshot{
		Symbol:     c.symbol,
		Qty:        c.qty,
		AvePrice:   c.avgPrice,
		Realized:   c.realized,
		Unrealized: c.unrealized,
		MarkPrice:  c.markPrice,
		UpdateAt:   c.updateAt,
	}
}
