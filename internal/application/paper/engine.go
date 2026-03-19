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
	"quant-trading/internal/application/risk"
	"quant-trading/internal/domain/execution"
	execution2 "quant-trading/internal/domain/order"
	"quant-trading/internal/infrastructure/broker"
)

// Engine 纸上交易引擎（与 backtest.Engine 对称）
type Engine struct {
	broker     broker.Broker
	accountCtx *account.Context
}

func (e *Engine) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (e *Engine) Evaluate() {
	//TODO implement me
	panic("implement me")
}

func (e *Engine) Results() <-chan *risk.Result {
	//TODO implement me
	panic("implement me")
}

func NewEngine(broker broker.Broker, accountCtx *account.Context) *Engine {
	return &Engine{
		broker:     broker,
		accountCtx: accountCtx,
	}
}

// SubmitOrder 供策略 / dispatcher 调用
func (e *Engine) SubmitOrder(ctx context.Context, ord *execution2.Order) error {
	_, err := e.broker.SubmitOrder(ctx, ord)
	return err
}

// Start 启动事件订阅（paper 模式下即时成交）
func (e *Engine) Start(ctx context.Context) error {
	// 订阅 Broker 事件 -> 更新 account
	go func() {
		for evt := range e.broker.SubscribeEvents(ctx) {
			if evt.Type == execution.OrderFilled {
				e.accountCtx.ApplyFill(evt.Symbol, evt.Side, evt.Price, evt.FilledQty)
			}
		}
	}()
	return nil
}
