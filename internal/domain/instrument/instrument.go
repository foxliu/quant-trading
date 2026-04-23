package instrument

import (
	"time"
)

/*
	Instrument  合约模型

股票：
Type = STOCK
Multiplier = 1
其余为空
期货：
Type = FUTURE
ExpiryDate + Multiplier
期权：
Type = OPTION
Strike + OptionType + Underlying
*/
type Instrument struct {
	ID       string
	Symbol   string // BTCUSDT / AAPL
	Exchange string // binance / nasdaq
	Type     Type

	// ===== 衍生品属性（股票为空）=====
	ExpiryDate *time.Time // 期货 / 期权  合约到期日
	Multiplier float64    // 合约乘数（股票=1）

	// ===== 期权专属 =====
	OptionType OptionType // CALL / PUT
	Strike     float64
	Underlying string // 标的 Symbol（如 510050)
}
