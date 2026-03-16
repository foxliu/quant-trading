package position

import (
	"errors"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/trade"
	"quant-trading/pkg/utils"
	"time"
)

func (c *Context) OnExecutionEvent(evt *execution.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch evt.Type {
	case execution.OrderFilled:
		return c.applyFill(evt)

	case execution.OrderPartiallyFilled:
		return c.applyFill(evt)

	default:
		// Accepted / StatusRejected / StatusCanceled 不影响 Position
		return nil
	}
}

func (c *Context) applyFill(evt *execution.Event) error {
	if evt.FilledQty <= 0 {
		return errors.New("无效的增持仓数")
	}

	now := evt.Timestamp
	if now.IsZero() {
		now = time.Now()
	}

	// Buy = +qty, Sell = -qty
	delta := evt.FilledQty
	if evt.Side == trade.Sell {
		delta = -delta
	}

	// === 无仓位 → 新开仓 ===
	if c.pos == nil {
		c.pos = &trade.Position{
			Symbol:   c.symbol,
			Qty:      delta,
			AvgPrice: evt.Price,
			OpenTime: now,
			UpdateAt: now,
		}
		return nil
	}

	prevQty := c.pos.Qty
	newQty := prevQty + delta

	// == 完全平仓 ===
	if newQty == 0 {
		c.pos = nil
		return nil
	}

	// === 同方向加仓 ===
	if (prevQty > 0 && delta > 0) || (prevQty < 0 && delta < 0) {
		totalCost := float64(utils.Abs(prevQty))*c.pos.AvgPrice + float64(utils.Abs(delta))*evt.Price

		c.pos.Qty = newQty
		c.pos.AvgPrice = totalCost / float64(utils.Abs(newQty))
		c.pos.UpdateAt = now
		return nil
	}

	// === 反方向成交 → 减仓 / 平仓 / 反手 ===
	if utils.Abs(delta) < utils.Abs(prevQty) {
		// 减仓，不变成本
		c.pos.Qty = newQty
		c.pos.UpdateAt = now
		return nil
	}

	// === 反手 ===
	c.pos.Qty = newQty
	c.pos.AvgPrice = evt.Price
	c.pos.OpenTime = now
	c.pos.UpdateAt = now

	return nil
}
