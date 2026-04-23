package risk

import "time"

// Status 风控整体状态
type Status struct {
	IsBlocked     bool      `json:"is_blocked"`
	BlockReason   string    `json:"block_reason"`
	LastCheckTime time.Time `json:"last_check_time"`
	EmergencyStop bool      `json:"emergency_stop"`
	ActiveRules   int       `json:"active_rules"`
}
