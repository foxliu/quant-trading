/*
Package paper
--------------------------------------------------------------------------------
CAUTION: TEMPORARY PACKAGE
Deprecated: 这个包是临时过渡使用的
*/
package paper

import (
	"context"
	"quant-trading/internal/application/account"
	aExecution "quant-trading/internal/application/execution"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/infrastructure/logger"
	"sync"

	"go.uber.org/zap"
)

// executionBackend 下单并推送执行事件（CTP / Paper 等适配器实现，独立于 broker.Broker 主接口）
type executionBackend interface {
	SubmitOrder(ctx context.Context, ord *order.Order) (string, error)
	SubscribeEvents(ctx context.Context) <-chan execution.Event
}

// Engine 纸上交易引擎（与 backtest.Engine 对称）
type Engine struct {
	broker      executionBackend
	accountCtx  *account.Context
	listeners   []aExecution.Listener
	listenersMu sync.RWMutex
}

func NewEngine(broker executionBackend, accountCtx *account.Context) *Engine {
	return &Engine{
		broker:     broker,
		accountCtx: accountCtx,
		listeners:  make([]aExecution.Listener, 0),
	}
}

// Submit 供策略 / dispatcher 调用
func (e *Engine) Submit(ctx context.Context, ord *order.Order) error {
	_, err := e.broker.SubmitOrder(ctx, ord)
	return err
}

func (e *Engine) RegisterListener(listener aExecution.Listener) {
	e.listenersMu.Lock()
	e.listeners = append(e.listeners, listener)
	e.listenersMu.Unlock()
}

// Start 启动事件订阅（paper 模式下即时成交）
func (e *Engine) Start(ctx context.Context) error {
	go func() {
		for evt := range e.broker.SubscribeEvents(ctx) {
			e.listenersMu.RLock()
			for _, l := range e.listeners {
				l.OnExecutionEvent(ctx, &evt)
			}
			e.listenersMu.RUnlock()
		}
	}()
	logger.Logger.With(zap.String("module", "paper.engine")).Info("Paper Engine 已启动")
	return nil
}
