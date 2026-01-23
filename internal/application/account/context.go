package account

import (
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/execution"
	"sync"
	"time"
)

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
	}
}
