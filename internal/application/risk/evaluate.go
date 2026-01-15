package risk

import "quant-trading/internal/domain/strategy"

// 风控 + 决策主流程（内部核心）

func (e *engine) evaluate(ctx Context, signal strategy.Signal) Decision {
	if !e.passRiskCheck(ctx, signal) {
		return Decision{}
	}

	return e.planOrders(ctx, signal)
}
