package pnl

import (
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/trade"
	"quant-trading/pkg/utils"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Revaluate(
	balance account.Balance,
	pos *trade.Position,
	marketPrice float64,
) Result {

	var unrealized float64
	var exposure float64

	if pos != nil {
		//direction := float64(1)
		//if pos.Qty < 0 {
		//	direction = -1
		//}
		unrealized = float64(pos.Qty) * (marketPrice - pos.AvgPrice)
		exposure = float64(utils.Abs(pos.Qty)) * marketPrice
	}

	equity := balance.Available() + unrealized

	return Result{
		UnrealizedPnL: unrealized,
		Exposure:      exposure,
		Equity:        equity,
	}
}
