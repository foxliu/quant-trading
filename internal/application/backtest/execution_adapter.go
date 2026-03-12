package backtest

import (
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/execution"
	execution2 "quant-trading/internal/domain/execution"
)

type ExecutionAdapter struct {
	account        *account.Context
	orderManager   *execution.OrderManager
	matchingEngine *execution.MatchingEngine
}

func NewExecutionAdapter(acc *account.Context) *ExecutionAdapter {
	return &ExecutionAdapter{
		account:        acc,
		orderManager:   execution.NewOrderManager(),
		matchingEngine: execution.NewMatchingEngine(0.0005),
	}
}

func (e *ExecutionAdapter) SubmitOrder(
	accountID string,
	strategyID string,
	symbol string,
	side execution2.Side,
	qty int64,
	marketPrice float64,
) error {
	createdOrder := e.orderManager.CreateOrder(accountID, strategyID, symbol, side, qty, marketPrice)

	//fill := e.matchingEngine.Match(createdOrder, marketPrice)

	e.account.ApplyFill(symbol, side, marketPrice, qty)
	e.orderManager.MarkFilled(createdOrder.ID())
	return nil
}
