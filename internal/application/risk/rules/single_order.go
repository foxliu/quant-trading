package rules

import (
	"fmt"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
)

// SingleOrderLimitRule 单笔订单最大数量限制
type SingleOrderLimitRule struct {
	maxQty int64
}

func NewSingleOrderLimitRule(maxQty int64) risk.Rule {
	return &SingleOrderLimitRule{
		maxQty: maxQty,
	}
}

func (r *SingleOrderLimitRule) Name() string {
	return "SingleOrderLimitRule"
}

func (r *SingleOrderLimitRule) Type() risk.RuleType {
	return risk.RuleSingleOrder
}

func (r *SingleOrderLimitRule) CheckOrder(ord *order.Order) risk.CheckResult {
	if ord.Qty() > r.maxQty {
		return risk.CheckResult{
			Action:   risk.ActionBlock,
			RuleType: r.Type(),
			Message:  fmt.Sprintf("单笔下单数量 %d 超过最大下单限制 %d", ord.Qty(), r.maxQty),
			Level:    1,
		}
	}
	return risk.CheckResult{Action: risk.ActionAllow}
}

func (r *SingleOrderLimitRule) CheckPosition() risk.CheckResult {
	return risk.CheckResult{Action: risk.ActionAllow}
}
