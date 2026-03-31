package risk

import (
	"quant-trading/internal/application/risk/rules"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/risk"
	"quant-trading/internal/infrastructure/logger"
	"sync"

	"go.uber.org/zap"
)

// Engine 风控引擎实现
type Engine struct {
	mu     sync.RWMutex
	rules  []risk.Rule
	status risk.Status
	logger *zap.Logger
}

func NewEngine() *Engine {
	e := &Engine{
		rules:  make([]risk.Rule, 0),
		logger: logger.Logger.With(zap.String("module", "risk.engine")),
	}

	// 添加默认规则
	e.RegisterRule(rules.NewMaxPositionRule(20))
	e.RegisterRule(rules.NewSingleOrderLimitRule(5))
	e.RegisterRule(rules.NewMaxDrawDownRule(0.1))
	e.RegisterRule(rules.NewDailyLossLimitRule(80000))

	e.logger.Info("风控引擎初始化完成", zap.Int("rule_count", len(e.rules)))
	return e
}

func (e *Engine) RegisterRule(rule risk.Rule) {
	e.mu.Lock()
	e.rules = append(e.rules, rule)
	e.mu.Unlock()
}

func (e *Engine) CheckOrder(ord *order.Order) risk.CheckResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rule := range e.rules {
		if result := rule.CheckOrder(ord); result.Action.IsBlock() {
			e.logger.Warn("风控拦截订单",
				zap.String("rule", result.RuleType.String()),
				zap.String("message", result.Message))
			return result
		}
	}
	return risk.CheckResult{Action: risk.ActionAllow}
}

func (e *Engine) CheckPosition() risk.CheckResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rule := range e.rules {
		if result := rule.CheckPosition(); result.Action.IsBlock() {
			e.logger.Warn("风控拦截持仓",
				zap.String("rule", result.RuleType.String()),
				zap.String("message", result.Message))
			return result
		}
	}
	return risk.CheckResult{Action: risk.ActionAllow}
}

func (e *Engine) GetStatus() risk.Status {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.status
}

func (e *Engine) EmergencyStop(reason string) {
	e.mu.Lock()
	e.status.EmergencyStop = true
	e.status.BlockReason = reason
	e.mu.Unlock()
	e.logger.Error("紧急熔断已触发", zap.String("reason", reason))
}
