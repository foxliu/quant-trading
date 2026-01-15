package risk

import "quant-trading/internal/domain/strategy"

func (e *engine) handleSignal(signal strategy.Signal) {
	ctx := e.ctxProvider.ContextFor(signal)

	decision := e.evaluate(ctx, signal)

	// TODO: 交给 Execution Engine
	_ = decision
}
