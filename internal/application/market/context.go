package market

import (
	"sync"
	"time"
)

type Price struct {
	Symbol    string
	Last      float64
	Timestamp time.Time
}

type Context struct {
	mu     sync.Mutex
	prices map[string]Price
}

func NewContext() *Context {
	return &Context{
		prices: make(map[string]Price),
	}
}

func (c *Context) Update(symbol string, price float64, ts time.Time) {
	if ts.IsZero() {
		ts = time.Now()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.prices[symbol] = Price{
		Symbol:    symbol,
		Last:      price,
		Timestamp: ts,
	}
}

func (c *Context) Get(symbol string) (Price, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	p, ok := c.prices[symbol]
	return p, ok
}
