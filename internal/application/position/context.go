package position

import (
	"errors"
	"quant-trading/internal/application/snapshot"
	"quant-trading/internal/domain/trade"
	"sync"
	"time"
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

func (c *Context) Take() snapshot.Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	var posCopy *trade.Position
	if c.pos != nil {
		p := *c.pos
		posCopy = &p
	}
	return &Snapshot{
		Symbol: c.symbol,
		Pos:    posCopy,
		At:     time.Now(),
	}
}

func (c *Context) Restore(s snapshot.Snapshot) error {
	ps, ok := s.(*Snapshot)
	if !ok {
		return errors.New("无效的 position snapshot")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.pos = ps.Pos
	return nil
}
