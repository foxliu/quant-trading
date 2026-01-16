package execution

import "time"

/*
EventType
=========

Execution Event 是“执行事实”，不是业务意图
*/
type EventType int

const (
	OrderAccepted EventType = iota
	OrderRejected
	OrderPartiallyFilled
	OrderFilled
	OrderCanceled
)

func (t EventType) String() string {
	switch t {
	case OrderAccepted:
		return "ORDER_ACCEPTED"
	case OrderRejected:
		return "ORDER_REJECTED"
	case OrderPartiallyFilled:
		return "ORDER_PARTIALLY_FILLED"
	case OrderFilled:
		return "ORDER_FILLED"
	case OrderCanceled:
		return "ORDER_CANCELED"
	default:
		return "UNKNOWN"
	}
}

/*
Event
=====

Execution Event 描述一次“执行侧发生的事实”
*/
type Event struct {
	OrderID string
	Type    EventType

	FilledQty int64   // 本次事件的成效数量（非累计）
	Price     float64 // 成交价（如适用）

	Timestamp time.Time
	Reason    string // 拒单 / 取消原因
}
