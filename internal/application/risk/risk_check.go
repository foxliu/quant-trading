package risk

import "quant-trading/internal/domain/strategy"

func (e *engine) passRiskCheck(ctx Context, signal strategy.Signal) bool {
	return true
}
