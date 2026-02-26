package event

import "time"

type Type string

const (
	// EventMarketPrice Market
	EventMarketPrice Type = "market.price"

	// EventSignal 策略产生的交易意图
	EventSignal Type = "strategy.signal"

	// EventOrderEvent Execution
	EventOrderEvent Type = "execution.order"

	// EventPositionChanged 仓位变更（可选）
	EventPositionChanged Type = "position.changed"

	// EventAccountUpdate 账户资金/权益更新（可选）
	EventAccountUpdate Type = "account.update"

	// EventRiskBreach 风控触发
	EventRiskBreach Type = "risk.breach"
)

/*
Envelope 是系统级的

Payload 是子系统级的
Event Bus 永远只认识 Envelope
*/
type Envelope struct {
	ID        uint64    `json:"id,omitempty"`
	Type      Type      `json:"type"`
	Source    string    `json:"source,omitempty"` // 如 "strategy-1"、"paper-engine"
	Timestamp time.Time `json:"timestamp"`

	Payload any `json:"payload"`
}
