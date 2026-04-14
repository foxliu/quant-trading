package execution

import (
	"quant-trading/internal/domain/order"
	"time"

	"github.com/pseudocodes/go2ctp/thost"
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

	Side     order.Side // Buy / Sell
	Quantity int64      // 本次事件的成效数量（非累计）
	Price    float64    // 成交价（如适用）

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

type AccountEvent struct {
	ReqID  int
	Data   *thost.CThostFtdcTradingAccountField
	IsLast bool
	Err    error
}
