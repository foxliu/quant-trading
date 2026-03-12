package backtest

import (
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/order"
)

type OrderExecutor struct {
	adapter *ExecutionAdapter
}

func NewOrderExecutor(adapter *ExecutionAdapter) *OrderExecutor {
	return &OrderExecutor{adapter: adapter}
}

// Handle 接收策略产生的订单事件
func (e *OrderExecutor) Handle(env *event.Envelope) {
	evt, ok := env.Payload.(*order.SubmitOrderEvent)
	if !ok {
		return
	}

	_ = e.adapter.SubmitOrder(
		evt.AccountID,
		evt.StrategyID,
		evt.Symbol,
		evt.Side,
		float64(evt.Quantity),
		evt.Price,
	)
}
