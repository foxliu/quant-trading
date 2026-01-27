package account

import (
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/trade"
	"sync"
	"time"
)

/*
Account 不计算盈亏
盈亏来自 Position Context
*/

type Context struct {
	mu sync.Mutex

	cfg     account.Config
	balance account.Balance

	updateAt time.Time
}

func NewContext(cfg account.Config) *Context {
	return &Context{
		cfg: cfg,
		balance: account.Balance{
			Cash:   cfg.InitialCash,
			Equity: cfg.InitialCash,
		},
		updateAt: time.Now(),
	}
}

func (c *Context) OnExecutionEvent(evt *execution.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch evt.Type {
	case execution.OrderFilled, execution.OrderPartiallyFilled:
		return c.applyFill(evt)

	case execution.FreeCharged:
		return c.applyFree(evt)
	default:
		return nil
	}
}

func (c *Context) applyFill(evt *execution.Event) error {
	cashDelta := float64(evt.FilledQty) * evt.Price

	// Buy: 扣现金
	if evt.Side == trade.Buy {
		c.balance.Cash -= cashDelta
	}

	// Sell: 回现金 + 已实现盈亏由 Position 决定
	if evt.Side == trade.Sell {
		c.balance.Cash += cashDelta
	}

	c.updateAt = evt.Timestamp
	return nil
}

func (c *Context) applyFree(evt *execution.Event) error {
	c.balance.Cash -= evt.Fee
	c.balance.RealizedPnL -= evt.Fee
	c.updateAt = evt.Timestamp
	return nil
}
