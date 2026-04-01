package rules

import (
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
)

/*
最大回撤（Floating Loss)
*/

// MaxDrawdownRule 最大回撤限制
type MaxDrawdownRule struct {
	maxDrawDown float64 // 例如 0.10 表示 10%
}

func NewMaxDrawDownRule(maxDrawDown float64) *MaxDrawdownRule {
	return &MaxDrawdownRule{
		maxDrawDown: maxDrawDown,
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
	// TODO: 实际项目中会从 accountCtx 获取当前权益和峰值权益
	// 这里简化返回 Allow，真实实现需注入 accountCtx
	return risk.CheckResult{Action: risk.ActionAllow}
}
