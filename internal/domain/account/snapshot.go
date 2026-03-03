package account

import (
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/portfolio"
)

type Snapshot struct {
	CapitalSnapshot   capital.Snapshot
	PortfolioSnapshot portfolio.Snapshot
}
