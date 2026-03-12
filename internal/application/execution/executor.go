package execution

import (
	"quant-trading/internal/domain/trade"
)

/*
真正下单
======
Execution 使用的是 trade.Side
qty 始终为正数
side 决定方向
*/

type Executor interface {
	MarketClose(symbol string, sid trade.Side, qty int64) error
}
