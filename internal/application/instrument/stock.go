package instrument

import (
	"errors"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/instrument"
)

/*
StockAdapter
============

股票适配器实现。

特点:
- 不需要保证金
- 最小交易单位通常为100股(1手)
- 不会到期
*/
type StockAdapter struct {
	BaseAdapter
	minLot int64 // 最小交易单位(股)
}

// NewStockAdapter 创建股票适配器
func NewStockAdapter() *StockAdapter {
	return &StockAdapter{
		BaseAdapter: BaseAdapter{
			instrumentType: instrument.Stock,
			multiplier:     1.0,
		},
		minLot: 100, // A股最小交易单位为100股
	}
}

// ValidateOrder 验证订单
func (a *StockAdapter) ValidateOrder(ord *execution.Order) error {
	// 检查数量是否为最小交易单位的整数倍
	if ord.Quantity%a.minLot != 0 {
		return errors.New("order quantity must be multiple of min lot size")
	}

	// 检查价格是否为正
	if ord.Price <= 0 {
		return errors.New("order price must be positive")
	}

	return nil
}

// CalculateMargin 股票不需要保证金
func (a *StockAdapter) CalculateMargin(qty int64, price float64) float64 {
	return 0
}

// CalculateValue 计算持仓价值
func (a *StockAdapter) CalculateValue(qty int64, price float64) float64 {
	return float64(qty) * price
}
