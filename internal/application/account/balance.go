package account

/*
未实现 PnL 在未来接入 Market Price Feed
*/

func (c *Context) recalculateEquity() {
	// TODO: 实现权益计算
	// equity := c.balance.Cash
	// for _, pos := range c.positions {
	//     equity += float64(pos.Qty) * pos.AvgPrice
	// }
	// c.balance.Equity = equity
	c.balance.Equity = c.TotalCapital()
}
