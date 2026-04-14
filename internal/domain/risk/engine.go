package risk

import (
	"quant-trading/internal/domain/order"
)

// RiskEngine 是 domain 层定义的接口（application 层实现）
type RiskEngine interface {
	CheckOrder(ord *order.Order) CheckResult
	CheckPosition() CheckResult
	GetStatus() Status
	EmergencyStop(reason string)
}
