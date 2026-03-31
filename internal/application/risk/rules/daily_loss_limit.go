package rules

import (
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
	"time"
)

// DailyLossLimitRule 日亏损限制
type DailyLossLimitRule struct {
	maxLoss    float64
	startOfDay time.Time
}

func NewDailyLossLimitRule(maxLoss float64) risk.Rule {
	return &DailyLossLimitRule{
		maxLoss:    maxLoss,
		startOfDay: time.Now().Truncate(24 * time.Hour),
	}
}

func (r *DailyLossLimitRule) Name() string {
	return "DailyLossLimit"
}

func (r *DailyLossLimitRule) Type() risk.RuleType {
	return risk.RuleDailyLossLimit
}

func (r *DailyLossLimitRule) CheckOrder(ord *order.Order) risk.CheckResult {
	// 日亏损检查通常在持仓/权益层面进行
	return risk.CheckResult{Action: risk.ActionAllow}
}

func (r *DailyLossLimitRule) CheckPosition() risk.CheckResult {
	// TODO: 实际项目中从 accountCtx 获取当日已实现亏损
	// 这里简化实现
	return risk.CheckResult{Action: risk.ActionAllow}
}
