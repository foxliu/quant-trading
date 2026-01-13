package risk

import "quant-trading/internal/domain/strategy"

// 仓位管理模块

type PositionManager interface {
	Apply(signal strategy.Signal) error
}
