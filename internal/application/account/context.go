package account

import (
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/portfolio"
	"time"
)

/*
Account 不计算盈亏
盈亏来自 Position Context
*/

/*
Context
=======

账户上下文。

职责:

1 管理账户状态
2 提供策略访问接口
3 聚合 capital / portfolio / balance
*/
type Context struct {
	acc *account.Account

	balance   *account.Balance
	capital   capital.Engine
	portfolio portfolio.Engine

	realizedPnL float64
}

func NewContext(acc *account.Account, cap capital.Engine, port portfolio.Engine) *Context {
	balance := account.NewBalance(cap.Available())
	return &Context{
		acc:       acc,
		balance:   balance,
		capital:   cap,
		portfolio: port,
	}
}

func (c *Context) AccountID() string {
	return c.acc.AccountID
}

func (c *Context) Available() float64 {
	return c.balance.Available()
}

func (c *Context) TotalCapital() float64 {
	return c.capital.Total()
}

func (c *Context) Equity() float64 {
	unrealized := c.portfolio.UnrealizedPnL()
	return c.capital.Total() + unrealized
}

func (c *Context) RealizedPnL() float64 {
	return c.realizedPnL
}

/*
Snapshot
========
生成账户快照
*/
func (c *Context) Snapshot() account.Snapshot {
	return account.Snapshot{
		AccountID:   c.acc.AccountID,
		Balance:     c.balance.Snapshot(),
		Capital:     c.capital.Snapshot(),
		Portfolio:   c.portfolio.Snapshot(),
		RealizedPnL: c.realizedPnL,
		Timestamp:   time.Now(),
	}
}

/*
Restore
=======
恢复账户状态

注意:
仅用于 replay / checkpoint
*/
func (c *Context) Restore(s account.Snapshot) {
	c.balance.Restore(s.Balance)
	c.capital.Restore(s.Capital)
	c.portfolio.Restore(s.Portfolio)
	c.realizedPnL = s.RealizedPnL
}

func (c *Context) ApplyFill(symbol string, side order.Side, price float64, qty int64) {
	c.portfolio.UpdateFill(symbol, side, price, qty)
}
