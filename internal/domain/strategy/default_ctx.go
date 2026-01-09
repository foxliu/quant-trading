package strategy

import (
	"quant-trading/internal/domain/market"
	"sync"
	"time"
)

/*
DefaultContext
==============

Context 的默认实现。

特点：
- 线程不安全（由 Runtime 串行保证）
- 简单 map 存储
- 可被替换为更高级实现（如指标缓存）
*/
type DefaultContext struct {
	now          time.Time
	currentEvent market.Event

	params map[string]interface{}
	state  map[string]interface{}

	mu sync.RWMutex
}

// NewContext
//
// 每个 Runtime 创建一个 Context 实例
func NewContext() Context {
	return &DefaultContext{
		params: make(map[string]any),
		state:  make(map[string]any),
	}
}

// ========= 时间 =========

func (c *DefaultContext) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.now
}

func (c *DefaultContext) SetNow(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// ========= 行情 =========

func (c *DefaultContext) CurrentEvent() market.Event {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentEvent
}

func (c *DefaultContext) SetCurrentEvent(event market.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentEvent = event
}

// ========= 私有状态 =========

func (c *DefaultContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state[key] = value
}

func (c *DefaultContext) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.state[key]
	return v, ok
}

func (c *DefaultContext) MustGet(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.state[key]
	if !ok {
		panic("strategy context key not found: " + key)
	}
	return v
}

func (c *DefaultContext) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.state, key)
}

// ========= 参数 =========

func (c *DefaultContext) Params() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.params
}

func (c *DefaultContext) SetParams(params map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.params = params
}
