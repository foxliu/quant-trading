package rules

import (
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
)

/*
最大回撤（Floating Loss)
*/

// MaxDrawDownRule 最大回撤限制
type MaxDrawDownRule struct {
	maxDrawDown float64 // 例如 0.10 表示 10%
}

func NewMaxDrawDownRule(maxDrawDown float64) *MaxDrawDownRule {
	return &MaxDrawDownRule{
		maxDrawDown: maxDrawDown,
	}
}

func (r *MaxDrawDownRule) Name() string {
	return "MaxDrawdown"
}

func (r *MaxDrawDownRule) Type() risk.RuleType {
	return risk.RuleMaxDrawdown
}

func (r *MaxDrawDownRule) CheckOrder(ord *order.Order) risk.CheckResult {
	// 下单前通常不检查回撤， 持仓检查时使用
	return risk.CheckResult{Action: risk.ActionAllow}
}

func (r *MaxDrawDownRule) CheckPosition() risk.CheckResult {
	// TODO: 实际项目中会从 accountCtx 获取当前权益和峰值权益
	// 这里简化返回 Allow，真实实现需注入 accountCtx
	return risk.CheckResult{Action: risk.ActionAllow}
}
