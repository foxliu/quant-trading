package instrument

import (
	"errors"
	"quant-trading/internal/domain/instrument"
	"sync"
)

/*
Context
=======

资产上下文负责管理所有资产适配器。

设计原则:
- 统一管理不同类型资产
- 提供资产查询接口
- 支持动态添加资产
*/
type Context struct {
	mu       sync.RWMutex
	adapters map[string]Adapter // symbol -> Adapter
}

// NewContext 创建资产上下文
func NewContext() *Context {
	return &Context{
		adapters: make(map[string]Adapter),
	}
}

// Register 注册资产适配器
func (c *Context) Register(symbol string, adapter Adapter) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.adapters[symbol]; exists {
		return errors.New("instrument already registered: " + symbol)
	}

	c.adapters[symbol] = adapter
	return nil
}

// Get 获取资产适配器
func (c *Context) Get(symbol string) (Adapter, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	adapter, exists := c.adapters[symbol]
	if !exists {
		return nil, errors.New("instrument not found: " + symbol)
	}

	return adapter, nil
}

// GetType 获取资产类型
func (c *Context) GetType(symbol string) (instrument.Type, error) {
	adapter, err := c.Get(symbol)
	if err != nil {
		return "", err
	}

	return adapter.GetType(), nil
}

// List 列出所有已注册的资产
func (c *Context) List() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	symbols := make([]string, 0, len(c.adapters))
	for symbol := range c.adapters {
		symbols = append(symbols, symbol)
	}

	return symbols
}

// Remove 移除资产适配器
func (c *Context) Remove(symbol string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.adapters[symbol]; !exists {
		return errors.New("instrument not found: " + symbol)
	}

	delete(c.adapters, symbol)
	return nil
}
