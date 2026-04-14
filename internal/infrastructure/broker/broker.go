package broker

import (
	"context"
	"quant-trading/internal/domain/order"
)

// Broker 定义基础设施层 Broker 接口（所有券商适配器必须实现）
// 1. 主接口（供 engine 依赖注入使用）
type Broker interface {
	Trader() TraderAPI // 交易相关操作
	MarketData() MdAPI // 行情订阅
	Connect(ctx context.Context) error
	Disconnect() error
	IsConnected() bool
	GetBroketType() string // "ctp" / "sim" (未来支持多券商时扩展）
}

// TraderAPI 交易相关操作
// 交易操作接口（TraderAPI）—— engine 最常用
type TraderAPI interface {
	SubmitOrder(ctx context.Context, ord *order.Order) (string, error)            // 下单
	CancelOrder(ctx context.Context, ord *order.Order) error                      // 撤单
	QueryOrderStatus(ctx context.Context, ord *order.Order) (*order.Order, error) // 查询订单状态
	// 禁止添加 GetBalance / GetPositions 等任何账户查询方法（由AccountSnapshot 和 EventBus 承担）
}

// MdAPI 行情订阅
// 行情订阅接口（MdAPI）—— 策略引擎 / Risk 引擎可能需要
type MdAPI interface {
	Subscribe(ctx context.Context, instruments []string) error
	Unsubscribe(instruments []string) error
	// 可选：AddTickHandler(handler func(Tick)) —— 如果需要直接回调行情
}

// internalAccountRefresher 账户刷新接口（内部使用）
// 内部专用接口（仅供 AccountRefreshScheduler 使用，unexported 或单独包）
type internalAccountRefresher interface {
	RefreshAccount(ctx context.Context, accountID string) error // 触发ReqQryTradingAccount + ReqQryInvestorPosition
}
