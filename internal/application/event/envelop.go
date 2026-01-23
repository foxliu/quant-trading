package event

import "time"

type Type string
type Source string

const (
	// EventMarketPrice Market
	EventMarketPrice Type = "market.price"

	// EventOrderEvent Execution
	EventOrderEvent Type = "execution.order"

	// EventPositionChanged Position
	EventPositionChanged Type = "position.changed"

	// EventRiskBreach Risk
	EventRiskBreach Type = "risk.breach"
)

/*
Envelope 是系统级的

Payload 是子系统级的
Event Bus 永远只认识 Envelope
*/
type Envelope struct {
	ID        uint64
	Type      Type
	Source    Source
	Timestamp time.Time

	Payload any
}
