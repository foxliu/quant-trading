package order

import (
	"time"
)

// OrderType 订单类型（未来扩展 Stop、StopLimit 等）
type OrderType int

const (
	TypeMarket OrderType = iota // 市价单
	TypeLimit                   // 限价单
)

/*
Order
=====

Order 表示一笔【可执行的标准化交易意图】。

约定：
- Order 是 Risk Engine 的最终产物
- Execution Engine 必须能够直接消费 Order
- Order 不表达策略意图（IntentLong / IntentShort）
- 只表达交易动作（Buy / Sell）
- Quantity 是相对变化量（Delta）
- Intent 复用 strategy.PositionIntent
- Side 使用 trade.Side
*/
type Order struct {
	// === 标识 ===
	id         string
	strategyID string
	accountID  string
	symbol     string
	side       Side
	orderType  OrderType
	qty        float64
	price      float64
	status     Status
	createdAt  time.Time
}

func NewOrder(id, strategyID, accountID, symbol string, side Side, orderType OrderType, price float64, qty float64) *Order {
	return &Order{
		id:         id,
		strategyID: strategyID,
		accountID:  accountID,
		symbol:     symbol,
		side:       side,
		orderType:  orderType,
		qty:        qty,
		price:      price,
		createdAt:  time.Now(),
	}
}

func (o *Order) ID() string           { return o.id }
func (o *Order) StrategyID() string   { return o.strategyID }
func (o *Order) AccountID() string    { return o.accountID }
func (o *Order) Symbol() string       { return o.symbol }
func (o *Order) Side() Side           { return o.side }
func (o *Order) OrderType() OrderType { return o.orderType }
func (o *Order) Qty() float64         { return o.qty }
func (o *Order) Price() float64       { return o.price }
func (o *Order) Status() Status       { return o.status }
func (o *Order) CreateAt() time.Time  { return o.createdAt }

func (o *Order) MarkFilled() {
	o.status = StatusFilled
}
