package market

import (
	"quant-trading/internal/domain/instrument"
	"time"
)

/*
Bar 行情模型

Domain 不关心：
主力合约
连续合约
那是 Data / Strategy 层的职责
*/
type Bar struct {
	Instrument instrument.Instrument
	Time       time.Time

	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}
