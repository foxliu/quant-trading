package execution

import (
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/execution"
)

type MatchingEngine struct {
	feeRate float64
}

func NewMatchingEngine(feeRate float64) *MatchingEngine {
	return &MatchingEngine{feeRate: feeRate}
}

// Match 回测简化： 全部按当前价格立即成交
func (m *MatchingEngine) Match(o *execution.Order, marketPrice float64) account.Fill {
	fee := marketPrice * o.Qty() * m.feeRate

	return account.Fill{
		Symbol: o.Symbol(),
		Side:   convertSide(o.Side()),
		Qty:    o.Qty(),
		Price:  marketPrice,
		Fee:    fee,
	}
}

func convertSide(s execution.Side) account.Side {
	if s == execution.Buy {
		return account.Buy
	}
	return account.Sell
}
