package event

import "time"

/*
Envelope 是系统级的

Payload 是子系统级的
Event Bus 永远只认识 Envelope
*/
type Envelope struct {
	//ID        uint64    `json:"id,omitempty"`
	Type      Type      `json:"type"`
	Source    string    `json:"source,omitempty"` // 如 "strategy-1"、"paper-engine"
	Timestamp time.Time `json:"timestamp"`

	Payload any `json:"payload"`
}
