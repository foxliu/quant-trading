package position

import (
	"quant-trading/internal/domain/trade"
	"sync"
)

/*
Context
=======

Position Context 维护账户维度的仓位事实。
*/
type Context struct {
	mu       sync.RWMutex
	position map[string]*trade.Position // key = symbol
}

func NewContext() *Context {
	return &Context{
		position: make(map[string]*trade.Position),
	}
}

func (c *Context) Get(symbol string) *trade.Position {
	c.mu.RLock()
	defer c.mu.Unlock()
	return c.position[symbol]
}

func (c *Context) Set(pos *trade.Position) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.position[pos.Symbol] = pos
}
