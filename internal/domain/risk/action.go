package risk

import "fmt"

// Action 风控执行动作（强类型，推荐使用）
type Action int

const (
	ActionAllow         Action = iota // 允许通过
	ActionBlock                       // 拦截该操作（拒单）
	ActionWarn                        // 仅警告，不拦截
	ActionForceClose                  // 强制平仓
	ActionHaltTrading                 // 暂停整个交易（紧急熔断）
	ActionEmergencyStop               // 紧急停止所有策略
)

func (a Action) String() string {
	switch a {
	case ActionAllow:
		return "ALLOW"
	case ActionBlock:
		return "BLOCK"
	case ActionWarn:
		return "WARN"
	case ActionForceClose:
		return "FORCE_CLOSE"
	case ActionHaltTrading:
		return "HALT_TRADING"
	case ActionEmergencyStop:
		return "EMERGENCY_STOP"
	default:
		return fmt.Sprintf("UNKNOWN_ACTION(%d)", a)
	}
}

// IsBlock 返回是否需要阻断操作
func (a Action) IsBlock() bool {
	return a == ActionBlock || a == ActionEmergencyStop || a == ActionHaltTrading || a == ActionForceClose
}
