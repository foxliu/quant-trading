package risk

import (
	"context"
	"quant-trading/internal/application/event"
	"quant-trading/internal/application/execution"
	risk2 "quant-trading/internal/domain/risk"
	"sync"
	"time"
)

type coordinator struct {
	exec     execution.Controller
	resultCh <-chan *Result
	bus      event.Bus

	mu sync.Mutex

	// 防止重复强平
	active map[string]bool

	stopCh chan struct{}
}

func NewCoordinator(exec execution.Controller, results <-chan *Result, bus event.Bus) Coordinator {
	return &coordinator{
		exec:     exec,
		resultCh: results,
		bus:      bus,
		active:   make(map[string]bool),
		stopCh:   make(chan struct{}),
	}
}

func (c *coordinator) Start(ctx context.Context) error {
	go c.runResultChannel(ctx)
	if c.bus != nil {
		c.bus.Subscribe(event.EventRiskBreach, c.handleEventBusRiskBreach)
	}
	return nil
}

func (c *coordinator) runResultChannel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		case res := <-c.resultCh:
			c.OnRiskResult(res)
		}
	}
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
		c.handleHaltTrading(res)

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

func (c *coordinator) handleEventBusRiskBreach(evt *event.Envelope) {
	if breach, ok := evt.Payload.(*risk2.Breach); ok {
		res := &Result{
			RuleName: breach.RuleName,
			Action:   risk2.ActionForceClose,
			Reason:   breach.Reason,
			Time:     time.Now(),
		}
		c.OnRiskResult(res)
	}
}

func (c *coordinator) handleHaltTrading(res *Result) {
	// TODO: 实现暂停交易逻辑
	// 当前阶段仅记录日志
	_ = res
}

func extractSymbol(res *Result) string {
	// 当前阶段：单 symbol 系统
	// 当前单 symbol 系统可直接返回默认，或从 res.Reason 解析
	// 未来可改为 res.Meta["symbol"]
	return "DEFAULT"
}
