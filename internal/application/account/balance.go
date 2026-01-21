package account

/*
未实现 PnL 在未来接入 Market Price Feed
*/

func (c *Context) recalculateEquity() {
	equity := c.balance.Cash

	// 当前阶段：不引入行情价格，仅使用仓位名义价值
	// 当前阶段 Equity ≈ Cash + Cost Basis
	for _, pos := range c.positions {
		equity += float64(pos.Qty) * pos.AvgPrice
	}

	c.balance.Equity = equity
}
