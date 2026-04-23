package instrument

import (
	"errors"
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/order"
	"time"
)

/*
OptionType 期权类型
*/
type OptionType string

const (
	OptionTypeCall OptionType = "CALL" // 看涨期权
	OptionTypePut  OptionType = "PUT"  // 看跌期权
)

/*
OptionsAdapter
==============

期权适配器实现。

特点:
- 需要保证金(卖方)
- 有合约乘数
- 有到期日
- 有行权价
- 分为看涨和看跌
*/
type OptionsAdapter struct {
	BaseAdapter
	optionType  OptionType // 期权类型
	strikePrice float64    // 行权价
	expiryDate  time.Time  // 到期日
	marginRate  float64    // 保证金比例(卖方)
}

// NewOptionsAdapter 创建期权适配器
func NewOptionsAdapter(
	optionType OptionType,
	strikePrice float64,
	multiplier float64,
	marginRate float64,
	expiryDate time.Time,
) *OptionsAdapter {
	return &OptionsAdapter{
		BaseAdapter: BaseAdapter{
			instrumentType: instrument.Option,
			multiplier:     multiplier,
		},
		optionType:  optionType,
		strikePrice: strikePrice,
		expiryDate:  expiryDate,
		marginRate:  marginRate,
	}
}

// ValidateOrder 验证订单
func (a *OptionsAdapter) ValidateOrder(ord *order.Order) error {
	// 检查合约是否到期
	if a.IsExpired() {
		return errors.New("options contract has expired")
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

// CalculateMargin 计算保证金(仅卖方需要)
func (a *OptionsAdapter) CalculateMargin(qty int64, price float64) float64 {
	// 期权卖方保证金计算(简化版)
	// 实际计算会更复杂,需要考虑标的价格、波动率等
	return float64(qty) * price * a.multiplier * a.marginRate
}

// CalculateValue 计算持仓价值
func (a *OptionsAdapter) CalculateValue(qty int64, price float64) float64 {
	// 期权持仓价值 = 数量 × 权利金 × 合约乘数
	return float64(qty) * price * a.multiplier
}

// IsExpired 检查是否到期
func (a *OptionsAdapter) IsExpired() bool {
	return time.Now().After(a.expiryDate)
}

// GetStrikePrice 获取行权价
func (a *OptionsAdapter) GetStrikePrice() float64 {
	return a.strikePrice
}

// GetOptionType 获取期权类型
func (a *OptionsAdapter) GetOptionType() OptionType {
	return a.optionType
}

// CalculateIntrinsicValue 计算内在价值
func (a *OptionsAdapter) CalculateIntrinsicValue(underlyingPrice float64) float64 {
	if a.optionType == OptionTypeCall {
		// 看涨期权内在价值 = Max(标的价格 - 行权价, 0)
		if underlyingPrice > a.strikePrice {
			return underlyingPrice - a.strikePrice
		}
		return 0
	} else {
		// 看跌期权内在价值 = Max(行权价 - 标的价格, 0)
		if a.strikePrice > underlyingPrice {
			return a.strikePrice - underlyingPrice
		}
		return 0
	}
}
