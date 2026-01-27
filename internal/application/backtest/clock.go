package backtest

import (
	"sync"
	"time"
)

/*
Clock
=====

回测时钟负责管理回测中的时间流逝。

设计原则:
- 回测时间是可控的,不依赖系统时间
- 支持时间快进
- 线程安全
*/
type Clock struct {
	mu  sync.RWMutex
	now time.Time
}

// NewClock 创建回测时钟
func NewClock(startTime time.Time) *Clock {
	return &Clock{
		now: startTime,
	}
}

// Now 获取当前回测时间
func (c *Clock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.now
}

// SetNow 设置当前回测时间
func (c *Clock) SetNow(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// Advance 推进回测时间
func (c *Clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
