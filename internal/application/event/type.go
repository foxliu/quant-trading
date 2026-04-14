package event

type Type int

const (
	EventMarketPrice          Type = iota // 市场行情事件（K线、Tick）
	EventSignal                           // 策略信号事件（供 RiskContext / Execution 使用）
	EventOrderSubmitted                   // 订单提交
	EventOrderFilled                      // 订单成交
	EventOrderCanceled                    // 订单撤单
	EventOrderRejected                    // 订单拒单
	EventDisconnected                     // 连接断开
	EventCTPTradingAccountRtn             // CTP请求查询资金账户响应
	EventCTPOrderRtn                      // CTP订单响应
)

func (t Type) String() string {
	switch t {
	case EventMarketPrice:
		return "market.price"
	case EventSignal:
		return "strategy.signal"
	case EventOrderSubmitted:
		return "order.submitted"
	case EventOrderFilled:
		return "order.filled"
	case EventOrderCanceled:
		return "order.canceled"
	case EventOrderRejected:
		return "order.rejected"
	case EventDisconnected:
		return "disconnected"
	default:
		return "unknown"
	}
}
