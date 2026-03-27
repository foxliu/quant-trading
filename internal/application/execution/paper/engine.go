package paper

import (
	"context"
	"errors"
	"quant-trading/internal/application/execution"
	dExecution "quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"time"
)

/*
PaperExecution
==============

PaperExecution 是一个最小 Execution Engine 实现：
- 不接交易所
- 不模拟盘口
- 用于联调 / 回测 / 冻结接口
*/
type PaperExecution struct {
	listener execution.Listener
}

func NewPaperExecution() *PaperExecution {
	return &PaperExecution{}
}

func (e *PaperExecution) RegisterListener(l execution.Listener) {
	e.listener = l
}

func (e *PaperExecution) Submit(ctx context.Context, ord *order.Order) error {
	if e.listener == nil {
		return errors.New("execution listener not registered")
	}

	now := time.Now()

	// 1. Order Accept
	e.listener.OnExecutionEvent(ctx, &dExecution.Event{
		OrderID:   ord.ID(),
		Symbol:    ord.Symbol(),
		Type:      dExecution.EventOrderAccepted,
		Side:      ord.Side(),
		Timestamp: now,
	})

	// 2. 直接全部成交（最简单模型）
	e.listener.OnExecutionEvent(ctx, &dExecution.Event{
		OrderID:   ord.ID(),
		Symbol:    ord.Symbol(),
		Type:      dExecution.EventOrderFilled,
		Side:      ord.Side(),
		Quantity:  ord.Qty(),
		Price:     ord.Price(),
		Timestamp: now.Add(1 * time.Millisecond),
	})
	return nil
}
