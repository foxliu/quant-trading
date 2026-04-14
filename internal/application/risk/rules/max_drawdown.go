package rules

import (
	"fmt"
	"quant-trading/internal/application/account"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
)

/*
最大回撤（Floating Loss)
*/

// MaxDrawdownRule 最大回撤限制
type MaxDrawdownRule struct {
	maxDrawDown float64 // 例如 0.10 表示 10%
	accountCtx  *account.Context
}

func NewMaxDrawDownRule(maxDrawDown float64, accountCtx *account.Context) *MaxDrawdownRule {
	return &MaxDrawdownRule{
		maxDrawDown: maxDrawDown,
		accountCtx:  accountCtx,
	}
}

func (r *MaxDrawdownRule) Name() string {
	return "MaxDrawdown"
}

func (r *MaxDrawdownRule) Type() risk.RuleType {
	return risk.RuleMaxDrawdown
}

func (r *MaxDrawdownRule) CheckOrder(ord *order.Order) risk.CheckResult {
	// 下单前通常不检查回撤， 持仓检查时使用
	return risk.CheckResult{Action: risk.ActionAllow}
}

func (r *MaxDrawdownRule) CheckPosition() risk.CheckResult {
	// 从 accountCtx中获取当前回撤
	positions, _ := r.accountCtx.GetPositions()
	for _, p := range positions {
		if p.UnrealizedPnL() < 0 {
			drawdown := -p.UnrealizedPnL() / (float64(p.Quantity) * p.OpenPrice)
			if drawdown > r.maxDrawDown {
				return risk.CheckResult{
					Action:   risk.ActionBlock,
					RuleType: r.Type(),
					Message:  fmt.Sprintf("品种 %s 当前回撤 %.2f 超过最大回撤限制 %.2f", p.Instrument.Symbol, drawdown*100, r.maxDrawDown*100),
					Level:    2,
				}
			}
		}
	}
	return risk.CheckResult{Action: risk.ActionAllow}
}
