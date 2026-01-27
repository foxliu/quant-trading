package backtest

import (
	"context"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"time"
)

/*
Simulator
=========

模拟撮合引擎负责在回测中模拟订单成交。

设计原则:
- 简化的撮合逻辑(不考虑盘口深度)
- 支持手续费和滑点
- 立即成交模型
*/
type Simulator struct {
	commission float64 // 手续费率
	slippage   float64 // 滑点(百分比)
}

// NewSimulator 创建模拟撮合引擎
func NewSimulator(commission, slippage float64) *Simulator {
	return &Simulator{
		commission: commission,
		slippage:   slippage,
	}
}

// Submit 提交订单
func (s *Simulator) Submit(ctx context.Context, ord *order.Order) ([]*execution.Event, error) {
	now := time.Now()

	events := make([]*execution.Event, 0, 3)

	// 1. 订单接受事件
	events = append(events, &execution.Event{
		OrderID:   ord.OrderID,
		Symbol:    ord.Symbol,
		Type:      execution.OrderAccepted,
		Side:      ord.Side,
		Timestamp: now,
	})

	// 2. 计算成交价格(考虑滑点)
	fillPrice := ord.Price
	if s.slippage > 0 {
		fillPrice = ord.Price * (1 + s.slippage)
	}

	// 3. 订单完全成交事件
	events = append(events, &execution.Event{
		OrderID:   ord.OrderID,
		Symbol:    ord.Symbol,
		Type:      execution.OrderFilled,
		Side:      ord.Side,
		FilledQty: ord.Quantity,
		Price:     fillPrice,
		Timestamp: now.Add(1 * time.Millisecond),
	})

	// 4. 手续费事件
	if s.commission > 0 {
		fee := float64(ord.Quantity) * fillPrice * s.commission
		events = append(events, &execution.Event{
			OrderID:   ord.OrderID,
			Symbol:    ord.Symbol,
			Type:      execution.FreeCharged,
			Fee:       fee,
			Timestamp: now.Add(2 * time.Millisecond),
		})
	}

	return events, nil
}
