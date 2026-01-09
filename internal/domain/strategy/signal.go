package strategy

import (
	"quant-trading/internal/domain/common"
	"quant-trading/internal/domain/instrument"
	"time"
)

/*
Signal 策略信号
注意：
期权多腿 ≠ Domain 信号
多腿在 Decision 层组合
*/
type Signal struct {
	Instrument instrument.Instrument
	Time       time.Time

	Side     common.Side
	Strength float64 // 0-1, 用于多种策略仲裁
	Reason   string  // 可选： 调试、复盘
}
