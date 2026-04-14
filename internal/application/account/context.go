package account

import (
	"errors"
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
		AccountID:  c.acc.AccountID,
		Balance:    c.balance.Snapshot(),
		UpdateTime: time.Now(),
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

// GetPositions 返回所有持仓
func (c *Context) GetPositions() ([]portfolio.Position, error) {
	return c.portfolio.GetPositions()
}

// GetPosition 返回单个品种持仓
func (c *Context) GetPosition(symbol string) (portfolio.Position, error) {
	pos, ok := c.portfolio.GetPosition(symbol)
	if !ok {
		return portfolio.Position{}, errors.New("无此持仓")
	}
	return pos, nil
}

// GetMaxDrawdown 获取最大回撤比例（0.0 ~ 1.0）
func (c *Context) GetMaxDrawdown() float64 {
	return c.portfolio.GetMaxDrawdown()
}

// GetDailyRealizedPnL 获取当日已实现盈亏
func (c *Context) GetDailyRealizedPnL() float64 {
	return c.portfolio.GetDailyRealizedPnL()
}
