// Package portfolio  持仓
package portfolio

import (
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/order"
)

// PositionSide 持仓方向 (Long / Short)
type PositionSide int

const (
	PositionSideLong  PositionSide = 1
	PositionSideShort PositionSide = -1
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
	Quantity int64

	OpenPrice float64
	LastPrice float64

	MarginUsed float64 // 期货保证金, 股票=0
}

// IsLong 是否多头
func (p *Position) IsLong() bool {
	return p.Side == PositionSideLong
}

// IsShort 是否空头
func (p *Position) IsShort() bool {
	return p.Side == PositionSideShort
}

// Value 当前持仓价值
func (p *Position) Value() float64 {
	return float64(p.Quantity) * p.LastPrice
}

// UnrealizedPnL 未实现盈亏(浮动盈亏）
func (p *Position) UnrealizedPnL() float64 {
	return float64(p.Quantity) * (p.LastPrice - p.OpenPrice)
}

// ToOrderSide 转换为 OrderSide （用于下单等场景）
func (p *Position) ToOrderSide() order.Side {
	if p.IsLong() {
		return order.Buy
	}
	return order.Sell
}
