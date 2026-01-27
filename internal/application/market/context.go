package market

import (
	"sync"
	"time"
)

type Context struct {
	mu       sync.RWMutex
	price    float64
	updateAt time.Time
}

func (c *Context) OnMarketPrice(price float64, at time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.price = price
	c.updateAt = at
}

func (c *Context) Latest() (float64, time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.price, c.updateAt
}
