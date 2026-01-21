package risk

import (
	risk2 "quant-trading/internal/domain/risk"
	"time"
)

type Result struct {
	RuleName string
	Action   risk2.Action
	Reason   string
	Time     time.Time
}
