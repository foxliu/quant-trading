package risk

import "quant-trading/internal/domain/strategy"

type NoopPositionManger struct{}

func NewNoopPositionManager() *NoopPositionManger {
	return &NoopPositionManger{}
}

func (pm NoopPositionManger) Apply(signal strategy.Signal) error {
	return nil
}
