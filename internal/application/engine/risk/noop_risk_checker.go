package risk

import "quant-trading/internal/domain/strategy"

type NoopRiskChecker struct{}

func NewNoopRiskChecker() *NoopRiskChecker {
	return &NoopRiskChecker{}
}

func (c *NoopRiskChecker) Check(signal strategy.Signal) error {
	return nil
}
