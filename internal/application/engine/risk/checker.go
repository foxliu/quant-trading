package risk

// 风控模块
import "quant-trading/internal/domain/strategy"

type RiskChecker interface {
	Check(signal strategy.Signal) error
}
