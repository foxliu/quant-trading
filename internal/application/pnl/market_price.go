package pnl

import "time"

func (c *Context) OnMarketPrice(price float64, ts time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.qty == 0 {
		c.unrealized = 0
		return
	}

	if ts.IsZero() {
		ts = time.Now()
	}

	c.markPrice = price
	c.unrealized = float64(c.qty) * (price - c.avgPrice)
	c.updateAt = ts
}
