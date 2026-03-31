package risk

import (
	"quant-trading/internal/domain/order"
	"time"
)

// Status 风控整体状态
type Status struct {
	IsBlocked     bool      `json:"is_blocked"`
	BlockReason   string    `json:"block_reason"`
	LastCheckTime time.Time `json:"last_check_time"`
	EmergencyStop bool      `json:"emergency_stop"`
	ActiveRules   int       `json:"active_rules"`
}

// RiskEngine 是 domain 层定义的接口（application 层实现）
type RiskEngine interface {
	CheckOrder(ord *order.Order) CheckResult
	CheckPosition() CheckResult
	GetStatus() Status
	EmergencyStop(reason string)
}
