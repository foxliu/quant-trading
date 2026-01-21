package execution

import (
	"errors"
	"quant-trading/internal/domain/trade"
	"quant-trading/pkg/utils"
)

type controller struct {
	orders    OrderManager
	positions PositionReader
	executor  Executor
}

func NewController(orders OrderManager, positions PositionReader, executor Executor) Controller {
	return &controller{
		orders:    orders,
		positions: positions,
		executor:  executor,
	}
}

func (c *controller) Execute(cmd Command) error {
	switch cmd.Type {
	case CommandForceClose:
		return c.forceClose(cmd)
	case CommandCancelAll:
		return c.orders.CancelAll(cmd.Symbol)
	default:
		return errors.New("不支持 command: " + string(cmd.Type))
	}
}

func (c *controller) forceClose(cmd Command) error {
	symbol := cmd.Symbol

	// 1. 取消所有挂单
	_ = c.orders.CancelAll(symbol)

	// 2. 读取当前仓位
	pos := c.positions.GetPosition(symbol)
	if pos == nil || pos.Qty == 0 {
		return nil
	}

	// 3. 计算平仓方向
	var side trade.Side
	qty := utils.Abs(pos.Qty)

	if pos.Qty > 0 {
		// 多头 → 卖出
		side = trade.Sell
	} else {
		// 空头 → 买入
		side = trade.Buy
	}

	// 4. 市价强平
	return c.executor.MarketClose(symbol, side, qty)
}
