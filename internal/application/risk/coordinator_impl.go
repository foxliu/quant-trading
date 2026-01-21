package risk

import (
	"context"
	"quant-trading/internal/application/execution"
	risk2 "quant-trading/internal/domain/risk"
	"sync"
	"time"
)

type coordinator struct {
	exec execution.Controller

	resultCh <-chan *Result

	mu sync.Mutex

	// 防止重复强平
	active map[string]bool
}

func NewCoordinator(exec execution.Controller, results <-chan *Result) Coordinator {
	return &coordinator{
		exec:     exec,
		resultCh: results,
		active:   make(map[string]bool),
	}
}

func (c *coordinator) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case res := <-c.resultCh:
				c.OnRiskResult(res)
			}
		}
	}()
	return nil
}

func (c *coordinator) Stop() error {
	return nil
}

func (c *coordinator) OnRiskResult(res *Result) {
	if res == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	switch res.Action {
	case risk2.ActionForceClose:
		c.handleForceClose(res)

	case risk2.ActionHaltTrading:
		c.handleHalt(res)

	default:
		// ActionRejectOrder 交给 Execution 层前置校验
	}
}

/*
幂等：同一个 symbol 同时只能有一个强平
异步执行，防止阻塞 Risk Loop
*/
func (c *coordinator) handleForceClose(res *Result) {
	symbol := extractSymbol(res)

	if c.active[symbol] {
		return
	}

	c.active[symbol] = true

	cmd := execution.Command{
		Type:   execution.CommandForceClose,
		Symbol: symbol,
		Reason: res.Reason,
		Time:   time.Now(),
	}

	go func() {
		defer func() {
			c.mu.Lock()
			delete(c.active, symbol)
			c.mu.Unlock()
		}()
		_ = c.exec.Execute(cmd)
	}()
}

func extractSymbol(res *Result) string {
	// 当前阶段：单 symbol 系统
	return "DEFAULT"
}
