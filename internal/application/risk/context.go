package risk

import (
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/position"
	"quant-trading/internal/domain/pnl"
	"sync"
)

/*
Context
=======

Risk Context 表示一组风控规则的集合。
*/
type Context struct {
	Mu sync.Mutex

	Account  account.Snapshot
	Position position.Snapshot
	PnL      pnl.Snapshot
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) UpdateAccount(s account.Snapshot) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Account = s
}

func (c *Context) UpdatePosition(s position.Snapshot) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Position = s
}

func (c *Context) UpdatePnL(s pnl.Snapshot) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.PnL = s
}
