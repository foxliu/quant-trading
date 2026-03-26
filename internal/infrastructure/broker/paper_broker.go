package broker

import (
	"context"
	"fmt"
	"quant-trading/internal/application/account"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/trade"
	"sync"
	"time"
)

// PaperBroker 纸上交易（模拟）Broker 实现
type PaperBroker struct {
	mu         sync.RWMutex
	accountCtx *account.Context
	positions  map[string]trade.Position
	orders     map[string]*order.Order
	events     chan execution.Event
	commission float64 // 手续费
	slippage   float64 // 滑点率
	nextID     int
}

// NewPaperBroker 创建纸上交易 Broker
// accCtx: 账户上下文
// commission: 手续费 (0.0003 = 0.03%)
// slippage: 滑点率 (0.0001 = 0.01%)
func NewPaperBroker(accCtx *account.Context, commission, slippage float64) *PaperBroker {
	b := &PaperBroker{
		accountCtx: accCtx,
		positions:  make(map[string]trade.Position),
		orders:     make(map[string]*order.Order),
		events:     make(chan execution.Event),
		commission: commission,
		slippage:   slippage,
		nextID:     100000,
	}
	go b.runEventLoop() // 后台模拟成交
	return b
}

// SubmitOrder 实现Broker 接口
func (b *PaperBroker) SubmitOrder(ctx context.Context, ord *order.Order) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	orderID := fmt.Sprintf("P_%d", b.nextID)
	b.nextID++

	newOrd := order.NewOrder(orderID, ord.StrategyID(), ord.AccountID(), ord.Symbol(),
		ord.Side(), ord.OrderType(), ord.Price(), ord.Qty())

	b.orders[orderID] = newOrd
	b.simulateFill(newOrd)
	return orderID, nil
}

// CancelOrder 实现Broker 接口
func (b *PaperBroker) CancelOrder(ctx context.Context, orderID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	ord, ok := b.orders[orderID]
	if !ok || ord.Status() == order.StatusFilled || ord.Status() == order.StatusCanceled {
		return fmt.Errorf("订单不存在或不可撤单: %s", orderID)
	}
	ord.MarkFilled()

	b.events <- execution.Event{
		OrderID:   orderID,
		Symbol:    ord.Symbol(),
		Type:      execution.EventOrderCanceled,
		Side:      ord.Side(),
		FilledQty: ord.Qty(),
		Price:     ord.Price(),
		Timestamp: ord.CreateAt(),
	}
	return nil
}

// GetPositions 实现 Broker 接口
func (b *PaperBroker) GetPositions(ctx context.Context) ([]trade.Position, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	pos := make([]trade.Position, 0, len(b.positions))
	for _, p := range b.positions {
		pos = append(pos, p)
	}
	return pos, nil
}

// GetBalance 实现 Broker 接口
func (b *PaperBroker) GetBalance(ctx context.Context) (cash float64, equity float64, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.accountCtx.Available(), b.accountCtx.Equity(), nil
}

// SubscribeEvents 实现 Broker 适配器接口
func (b *PaperBroker) SubscribeEvents(ctx context.Context) <-chan execution.Event {
	return b.events
}

// simulateFill 内部模拟即时成交（支持市价/限价）
func (b *PaperBroker) simulateFill(ord *order.Order) {
	// 模拟滑点
	fillPrice := ord.Price()
	if ord.OrderType() == order.TypeMarket {
		fillPrice += fillPrice * b.slippage // 简单滑点
	}

	fillQty := ord.Qty()
	// 手续费计算（后续可移到 account portfolio 内部）
	//fee := fillPrice * fillQty * b.commission

	// 更新账户
	b.accountCtx.ApplyFill(ord.Symbol(), ord.Side(), fillPrice, fillQty)

	// 推送事件（供 execution engine、position、pnl 使用）
	b.events <- execution.Event{
		Type:      execution.EventOrderFilled,
		OrderID:   ord.ID(),
		Price:     fillPrice,
		FilledQty: fillQty,
		Timestamp: time.Now(),
	}
}

// runEventLoop 后台事件循环（可扩展为延迟成交）
func (b *PaperBroker) runEventLoop() {
	// 当前版本即时成交，未来可加入延迟/部分成交逻辑
}
