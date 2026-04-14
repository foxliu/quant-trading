package rules

import (
	"fmt"
	"quant-trading/internal/application/account"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
	"sync"
)

// DailyLossLimitRule 日亏损限制
type DailyLossLimitRule struct {
	maxLoss           float64
	accountCtx        *account.Context
	dailyPnL          float64
	currentTradingDay string
	mu                sync.RWMutex
}

func NewDailyLossLimitRule(maxLoss float64, accountCtx *account.Context) risk.Rule {
	return &DailyLossLimitRule{
		maxLoss:    maxLoss,
		accountCtx: accountCtx,
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
	r.mu.Lock()
	defer r.mu.Unlock()

	// 获取当前交易日（从CTP或accountCtx 获取）
	tradingDay := r.accountCtx.GetCurrentTradingDay()

	// 如果交易日切换，重置当日亏损
	if tradingDay != r.currentTradingDay {
		r.currentTradingDay = tradingDay
		r.dailyPnL = 0
	}

	dailyLoss := r.accountCtx.GetDailyRealizedPnL()

	if dailyLoss < -r.maxLoss {
		return risk.CheckResult{
			Action:   risk.ActionBlock,
			RuleType: r.Type(),
			Message:  fmt.Sprintf("当日已达到亏损 %.2f 元 已超过最大允许 %.2f 元", dailyLoss, r.maxLoss),
			Level:    2,
		}
	}

	return risk.CheckResult{Action: risk.ActionAllow}
}
