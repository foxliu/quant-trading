package broker

import (
	"context"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/trade"
)

// Broker 定义基础设施层 Broker 接口（所有券商适配器必须实现）
type Broker interface {
	// SubmitOrder 提交订单（市价/限价），返回订单ID
	SubmitOrder(ctx context.Context, ord *order.Order) (string, error)
	// CancelOrder 撤单
	CancelOrder(ctx context.Context, orderID string) error
	// GetPositions 获取当前持仓
	GetPositions(ctx context.Context) ([]trade.Position, error)
	// GetBalance 获取账户余额 (现金+权益）
	GetBalance(ctx context.Context) (float64, float64, error)
	// SubscribeEvents 订阅成交/订单事件（供execution engine 使用)
	SubscribeEvents(ctx context.Context) <-chan execution.Event
}
