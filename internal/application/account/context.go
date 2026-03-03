package account

import (
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/portfolio"
	"quant-trading/internal/domain/trade"
	"sync"
)

/*
Account 不计算盈亏
盈亏来自 Position Context
*/

type Context struct {
	mu sync.RWMutex

	acc account.Account

	capital   capital.Engine
	portfolio portfolio.Engine

	dirty bool
}

func NewContext(acc account.Account, cap capital.Engine, port portfolio.Engine) *Context {
	return &Context{
		acc:       acc,
		capital:   cap,
		portfolio: port,
	}
}

// AcceptOrder 订单接受
func (c *Context) AcceptOrder(o *order.Order) error {
	err := c.capital.Freeze(o.OrderID, o.Symbol, o.Price, float64(o.Quantity), o.Side)
	if err != nil {
		return err
	}
	c.makeDirty()
	return nil
}

func (c *Context) ApplyFill(orderID string, symbol string, side trade.Side, qty int64, price float64) error {
	amount := float64(qty) * price
	if err := c.capital.Commit(orderID, amount); err != nil {
		return err
	}

	c.portfolio.UpdateFill(symbol, side, qty, price)

	c.makeDirty()
	return nil
}

// ApplyCancel 订单接受
func (c *Context) ApplyCancel(orderID string) error {
	if err := c.capital.Release(orderID); err != nil {
		return err
	}
	c.makeDirty()
	return nil
}

// ApplyMarketPrice 行情更新
func (c *Context) ApplyMarketPrice(symbol string, price float64) {
	c.portfolio.UpdateMarkPrice(symbol, price)
	c.makeDirty()
}

func (c *Context) AccountID() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.acc.AccountID
}

func (c *Context) Available() float64 {
	return c.capital.Available()
}

func (c *Context) Frozen() float64 {
	return c.capital.Frozen()
}

func (c *Context) TotalCapital() float64 {
	return c.capital.Total()
}

func (c *Context) Equity() float64 {
	return c.capital.Total() + c.portfolio.UnrealizedPnL()
}

func (c *Context) RealizedPnL() float64 {
	return c.portfolio.RealizedPnL()
}

func (c *Context) UnrealizedPnL() float64 {
	return c.portfolio.UnrealizedPnL()
}

// ========= Snapshot支持 =========

func (c *Context) Snapshot() account.Snapshot {
	return account.Snapshot{
		CapitalSnapshot:   c.capital.Snapshot(),
		PortfolioSnapshot: c.portfolio.Snapshot(),
	}
}

func (c *Context) Restore(s account.Snapshot) {
	c.capital.Restore(s.CapitalSnapshot)
	c.portfolio.Restore(s.PortfolioSnapshot)
}

// ========= 冻结相关的接口 ========

func (c *Context) FreezeOrder(o *order.Order) error {
	err := c.capital.Freeze(
		o.OrderID,
		o.Symbol,
		o.Price,
		float64(o.Quantity),
		o.Side,
	)
	if err != nil {
		return err
	}
	c.makeDirty()
	return nil
}

func (c *Context) CommitOrder(orderID string, amount float64) error {
	err := c.capital.Commit(orderID, amount)
	if err != nil {
		return err
	}
	c.makeDirty()
	return nil
}

func (c *Context) ReleaseOrder(orderID string) error {
	err := c.capital.Release(orderID)
	if err != nil {
		return err
	}

	c.makeDirty()
	return nil
}

// ========= Dirty标记机制 =========

func (c *Context) makeDirty() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dirty = true
}

func (c *Context) ConsumeDirty() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.dirty {
		c.dirty = false
		return true
	}
	return false
}
