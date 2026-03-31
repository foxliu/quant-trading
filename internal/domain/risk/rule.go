package risk

import "quant-trading/internal/domain/order"

// RuleType 风控规则的类型
type RuleType string

const (
	RuleMaxPosition    RuleType = "max_position"
	RuleMaxDrawdown    RuleType = "max_drawdown"
	RuleSingleOrder    RuleType = "single_order"
	RuleCircuitBreaker RuleType = "circuit_breaker"
	RuleDailyLossLimit RuleType = "daily_loss_limit"
	RuleMaxOrderFreq   RuleType = "max_order_frequency"
)

type CheckResult struct {
	Action   Action
	RuleType RuleType
	Message  string
	Level    int // 0 = 低 1 = 中 2 = 高
}

// String 方法保持与其他领域类型一致
func (r RuleType) String() string {
	return string(r)
}

// Rule 风控规则接口
type Rule interface {
	Name() string
	Type() RuleType
	CheckOrder(ord *order.Order) CheckResult
	CheckPosition() CheckResult
}
