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
	mu sync.Mutex

	symbol string
	pos    *trade.Position
}

func NewContext(symbol string) *Context {
	return &Context{
		symbol: symbol,
	}
}

func (c *Context) Symbol() string {
	return c.symbol
}

func (c *Context) Position() *trade.Position {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pos == nil {
		return nil
	}

	// 返回拷贝，防止外部篡改
	p := *c.pos
	return &p
}
