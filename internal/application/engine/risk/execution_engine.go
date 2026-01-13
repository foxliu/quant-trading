package risk

import "quant-trading/internal/domain/strategy"

// 下单执行模块

type ExecutionEngine interface {
	Execute(signal strategy.Signal) error
}
