package account

import (
	"quant-trading/internal/application/execution"
	"quant-trading/internal/application/position"
)

/*
Account 只接收“已经确定的事实”
=====
Account 不直接参与成交计算，避免重复状态源。
*/

func (c *Context) OnPositionSnapshot(snapshot *position.Snapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if snapshot != nil {
		if snapshot.Qty == 0 {
			delete(c.positions, snapshot.Symbol)
		} else {
			c.positions[snapshot.Symbol] = snapshot
		}
	}
	c.recalculateEquity()
}

func (c *Context) ONExecutionEvent(evt *execution.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 当前阶段： Account 不直接处理成交细节
	// 真实的PnL 结算由 Position + PriceFree 推导
}
