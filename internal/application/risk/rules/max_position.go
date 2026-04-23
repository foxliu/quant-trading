package rules

import (
	"fmt"
	"quant-trading/internal/application/account"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
)

/*
最大仓位限制
*/

// MaxPositionRule 单品种最大仓位限制
type MaxPositionRule struct {
	maxTotal   int64
	accountCtx *account.Context
}

func NewMaxPositionRule(maxTotal int64, accountCtx *account.Context) risk.Rule {
	return &MaxPositionRule{
		maxTotal:   maxTotal,
		accountCtx: accountCtx,
	}
}

func (r *MaxPositionRule) Name() string {
	return "MaxPositionRule"
}

func (r *MaxPositionRule) Type() risk.RuleType {
	return risk.RuleMaxPosition
}

func (r *MaxPositionRule) CheckOrder(ord *order.Order) risk.CheckResult {
	// 当前权做笔下单检查（持仓检查在 CheckPosition 中实现）
	if ord.Qty() > r.maxTotal {
		return risk.CheckResult{
			Action:   risk.ActionBlock,
			RuleType: r.Type(),
			Message:  fmt.Sprintf("单笔下单数量 %d 超过最大持仓限制 %d", ord.Qty(), r.maxTotal),
			Level:    2,
		}
	}
	return risk.CheckResult{Action: risk.ActionAllow}
}

func (r *MaxPositionRule) CheckPosition() risk.CheckResult {
	// 持仓检查逻辑由 Engine 统一调用，此处可扩展
	positions, _ := r.accountCtx.GetPositions()
	for _, p := range positions {
		if p.Quantity > r.maxTotal {
			return risk.CheckResult{
				Action:   risk.ActionBlock,
				RuleType: r.Type(),
				Message:  fmt.Sprintf("持仓数量 %d 超过最大持仓限制 %d", p.Quantity, r.maxTotal),
				Level:    2,
			}
		}
	}
	return risk.CheckResult{Action: risk.ActionAllow}
}
