package market

import "time"

type EventType string

const (
	EventMarket EventType = "MARKET"
	EventSignal EventType = "SIGNAL"
	EventOrder  EventType = "ORDER"
	EventTrade  EventType = "TRADE"
)

type Event struct {
	Type EventType
	Time time.Time
	Data interface{}
}
