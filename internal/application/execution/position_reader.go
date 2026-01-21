package execution

import "quant-trading/internal/domain/trade"

/*
 - Execution 不能直接持有 Position Context
 - 只能通过 Reader 读快照
*/

type PositionReader interface {
	GetPosition(symbol string) *trade.Position
}
