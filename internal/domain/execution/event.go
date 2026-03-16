package execution

import (
	"quant-trading/internal/domain/order"
	"time"
)

/*
EventType
=========

Execution Event 是“执行事实”，不是业务意图
*/
type EventType string

const (
	OrderSubmitted       EventType = "SUBMITTED"
	OrderAccepted        EventType = "ACCEPTED"
	OrderPartiallyFilled EventType = "PARTIALLY_FILLED"
	OrderFilled          EventType = "FILLED"
	OrderCanceled        EventType = "CANCELED"
	OrderRejected        EventType = "REJECTED"

	FreeCharged EventType = "FREE_CHARGED"
)

/*
Event
=====

Execution Event 描述一次“执行侧发生的事实”
*/
type Event struct {
	OrderID string
	Symbol  string

	Type EventType

	Side      order.Side // Buy / Sell
	FilledQty float64    // 本次事件的成效数量（非累计）
	Price     float64    // 成交价（如适用）

	Timestamp time.Time
	Reason    string // 拒单 / 取消原因

	Fee float64 // 费用
}

type SubmitOrderEvent struct {
	AccountID  string
	StrategyID string
	Symbol     string
	Side       order.Side
	Quantity   float64
	Price      float64
}
