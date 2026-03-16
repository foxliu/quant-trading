package instrument

import (
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/order"
)

/*
Adapter
=======

资产适配器负责处理不同资产类型的差异。

设计原则:
- 策略对资产类型无感知
- 差异化逻辑封装在适配器中
- 支持股票/期货/期权
*/
type Adapter interface {
	// GetType 获取资产类型
	GetType() instrument.Type

	// ValidateOrder 验证订单是否符合资产规则
	ValidateOrder(ord *order.Order) error

	// CalculateMargin 计算保证金(期货/期权)
	CalculateMargin(qty int64, price float64) float64

	// CalculateValue 计算持仓价值
	CalculateValue(qty int64, price float64) float64

	// IsExpired 检查合约是否到期(期货/期权)
	IsExpired() bool

	// GetMultiplier 获取合约乘数(期货/期权)
	GetMultiplier() float64
}

/*
BaseAdapter
===========

基础适配器实现,提供通用功能。
*/
type BaseAdapter struct {
	instrumentType instrument.Type
	multiplier     float64
}

// GetType 获取资产类型
func (a *BaseAdapter) GetType() instrument.Type {
	return a.instrumentType
}

// GetMultiplier 获取合约乘数
func (a *BaseAdapter) GetMultiplier() float64 {
	return a.multiplier
}

// IsExpired 默认不会到期
func (a *BaseAdapter) IsExpired() bool {
	return false
}
