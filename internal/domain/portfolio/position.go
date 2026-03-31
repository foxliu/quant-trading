// Package portfolio  持仓
package portfolio

import (
	"quant-trading/internal/domain/instrument"
)

// Position 持仓模型
// 同时支持：
// 股票多头
// 期货多 / 空
// 期权多 / 空
// 不做盈亏计算（那是 RiskContext 层）
type Position struct {
	Instrument instrument.Instrument

	Side     PositionSide // Long / Short
	Quantity float64

	OpenPrice float64
	LastPrice float64

	MarginUsed float64 // 股票=0
}
