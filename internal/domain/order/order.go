package order

import (
	"quant-trading/internal/domain/common"
	"quant-trading/internal/domain/instrument"
	"time"
)

/*
Order 订单模型

Quantity 含义：
股票：股数
期货 / 期权：合约数
*/
type Order struct {
	ID         string
	Instrument instrument.Instrument

	Side     common.Side
	Type     Type
	Price    float64
	Quantity float64

	Status     Status
	CreateTime time.Time
	UpdateTime time.Time
}
