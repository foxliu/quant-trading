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

// Handle 接收风控/协调器链路中的订单提交事件（与 Risk Coordinator 载荷一致）
func (e *OrderExecutor) Handle(env *event.Envelope) {
	ord, ok := env.Payload.(*order.Order)
	if !ok {
		return
	}

	_ = e.adapter.SubmitOrder(
		ord.AccountID(),
		ord.StrategyID(),
		ord.Symbol(),
		ord.Side(),
		ord.OrderType(),
		ord.Qty(),
		ord.Price(),
	)
}
