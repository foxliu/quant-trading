package risk

import "quant-trading/internal/domain/strategy"

type NoopExecutionEngine struct{}

func NewNoopExecutionEngine() *NoopExecutionEngine {
	return &NoopExecutionEngine{}
}

func (e *NoopExecutionEngine) Execute(signal strategy.Signal) error {
	return nil
}
