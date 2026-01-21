package execution

import (
	"fmt"
	"quant-trading/internal/domain/trade"
)

/*
最小Executor实现
*/

type dummyExecutor struct{}

func NewDummyExecutor() Executor {
	return &dummyExecutor{}
}

func (e *dummyExecutor) MarketClose(symbol string, side trade.Side, qty int64) error {
	fmt.Printf("[Exec] market close %s %s qty=%d\n", symbol, side, qty)
	return nil
}
