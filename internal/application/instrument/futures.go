package instrument

import (
	"errors"
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/order"
	"time"
)

/*
FuturesAdapter
==============

期货适配器实现。

特点:
- 需要保证金
- 有合约乘数
- 有到期日
- 支持双向交易
*/
type FuturesAdapter struct {
	BaseAdapter
	marginRate float64   // 保证金比例
	expiryDate time.Time // 到期日
}

// NewFuturesAdapter 创建期货适配器
func NewFuturesAdapter(multiplier float64, marginRate float64, expiryDate time.Time) *FuturesAdapter {
	return &FuturesAdapter{
		BaseAdapter: BaseAdapter{
			instrumentType: instrument.Future,
			multiplier:     multiplier,
		},
		marginRate: marginRate,
		expiryDate: expiryDate,
	}
}

// ValidateOrder 验证订单
func (a *FuturesAdapter) ValidateOrder(ord *order.Order) error {
	// 检查合约是否到期
	if a.IsExpired() {
		return errors.New("futures contract has expired")
	}

	// 检查价格是否为正
	if ord.Price() <= 0 {
		return errors.New("order price must be positive")
	}

	// 检查数量是否为正
	if ord.Qty() <= 0 {
		return errors.New("order quantity must be positive")
	}

	return nil
}

// CalculateMargin 计算保证金
func (a *FuturesAdapter) CalculateMargin(qty int64, price float64) float64 {
	// 保证金 = 数量 × 价格 × 合约乘数 × 保证金比例
	return float64(qty) * price * a.multiplier * a.marginRate
}

// CalculateValue 计算持仓价值
func (a *FuturesAdapter) CalculateValue(qty int64, price float64) float64 {
	// 期货持仓价值 = 数量 × 价格 × 合约乘数
	return float64(qty) * price * a.multiplier
}

// IsExpired 检查是否到期
func (a *FuturesAdapter) IsExpired() bool {
	return time.Now().After(a.expiryDate)
}
