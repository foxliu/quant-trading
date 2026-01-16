package position

import (
	"errors"
	"quant-trading/internal/application/execution"
	"quant-trading/internal/domain/trade"
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
		// Accepted / Rejected / Canceled 不影响 Position
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

	side := trade.SideFromOrderID(evt.OrderID) // 外部约定

	// === 无仓位 → 新开仓 ===
	if c.pos == nil {
		c.pos = &trade.Position{
			Symbol:   c.symbol,
			Sid:      side,
			Qty:      evt.FilledQty,
			AvgPrice: evt.Price,
			OpenTime: now,
			UpdateAt: now,
		}
		return nil
	}

	// === 同方向加仓 ===
	if c.pos.Sid == side {
		totalCost := float64(c.pos.Qty)*c.pos.AvgPrice + float64(evt.FilledQty)*evt.Price

		c.pos.Qty += evt.FilledQty
		c.pos.AvgPrice = totalCost / float64(c.pos.Qty)
		c.pos.UpdateAt = now
		return nil
	}

	// === 反方向成交 → 减仓 / 平仓 / 反手 ===
	if evt.FilledQty < c.pos.Qty {
		c.pos.Qty -= evt.FilledQty
		c.pos.UpdateAt = now
		return nil
	}

	if evt.FilledQty == c.pos.Qty {
		// 完全平仓
		c.pos = nil
		return nil
	}

	// === 反手 ===
	newQty := evt.FilledQty - c.pos.Qty

	c.pos = &trade.Position{
		Symbol:   c.symbol,
		Sid:      side,
		Qty:      newQty,
		AvgPrice: evt.Price,
		OpenTime: now,
		UpdateAt: now,
	}
	return nil
}
