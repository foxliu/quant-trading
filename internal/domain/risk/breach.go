package risk

import (
	"quant-trading/pkg/utils"
	"time"
)

// Breach 表示一次风险违规触发（由 RiskContext Rule Engine 产生）
type Breach struct {
	// 基础标识
	ID       string `json:"id"`       // 全局唯一 ID（推荐 UUID 或 seq）
	RuleName string `json:"ruleName"` // 触发规则名称（如 "MaxDrawdown"、"MaxPosition"）
	Symbol   string `json:"symbol"`   // 关联标的（空表示账户级风控）

	// 违规细节
	Action   Action `json:"action"`   // 建议执行动作（ForceClose / HaltTrading / RejectOrder）
	Reason   string `json:"reason"`   // 详细原因（人类可读）
	Severity string `json:"severity"` // 严重程度（Info / Warning / Critical）

	// 时间与上下文
	Timestamp time.Time         `json:"timestamp"`
	Meta      map[string]string `json:"meta,omitempty"` // 扩展字段（如 CurrentEquity、CurrentQty）
}

func NewBreach(ruleName, symbol, reason string, action Action, severity string) *Breach {
	return &Breach{
		ID:        utils.GenerateID(), // 可使用 pkg/utils 或 snowflake
		RuleName:  ruleName,
		Symbol:    symbol,
		Action:    action,
		Reason:    reason,
		Severity:  severity,
		Timestamp: time.Now(),
		Meta:      make(map[string]string),
	}
}
