package examples

import (
	"quant-trading/internal/domain/market"
	"quant-trading/internal/domain/strategy"
)

/*
MACrossStrategy
===============

双均线交叉策略示例。

策略逻辑:
- 短期均线上穿长期均线时买入
- 短期均线下穿长期均线时卖出

参数:
- shortPeriod: 短期均线周期
- longPeriod: 长期均线周期
*/
type MACrossStrategy struct {
	name string

	// 参数
	shortPeriod int
	longPeriod  int

	// 状态
	shortMA []float64 // 短期均线价格序列
	longMA  []float64 // 长期均线价格序列
}

// NewMACrossStrategy 创建双均线交叉策略
func NewMACrossStrategy(name string, shortPeriod, longPeriod int) *MACrossStrategy {
	return &MACrossStrategy{
		name:        name,
		shortPeriod: shortPeriod,
		longPeriod:  longPeriod,
		shortMA:     make([]float64, 0, longPeriod),
		longMA:      make([]float64, 0, longPeriod),
	}
}

// Name 返回策略名称
func (s *MACrossStrategy) Name() string {
	return s.name
}

// OnInit 初始化策略
func (s *MACrossStrategy) OnInit(ctx strategy.Context) error {
	// 初始化策略状态
	s.shortMA = make([]float64, 0, s.longPeriod)
	s.longMA = make([]float64, 0, s.longPeriod)
	return nil
}

// OnMarketEvent 处理行情事件
func (s *MACrossStrategy) OnMarketEvent(ctx strategy.Context, evt market.Event) ([]strategy.Signal, error) {
	// 提取价格(简化处理,实际应该根据事件类型解析)
	price := s.extractPrice(evt)
	if price <= 0 {
		return nil, nil
	}

	// 更新价格序列
	s.shortMA = append(s.shortMA, price)
	s.longMA = append(s.longMA, price)

	// 保持序列长度
	if len(s.shortMA) > s.longPeriod {
		s.shortMA = s.shortMA[1:]
		s.longMA = s.longMA[1:]
	}

	// 数据不足,不产生信号
	if len(s.shortMA) < s.longPeriod {
		return nil, nil
	}

	// 计算均线
	shortAvg := s.calculateMA(s.shortMA, s.shortPeriod)
	longAvg := s.calculateMA(s.longMA, s.longPeriod)

	// 计算前一个周期的均线(用于判断交叉)
	prevShortAvg := s.calculateMA(s.shortMA[:len(s.shortMA)-1], s.shortPeriod)
	prevLongAvg := s.calculateMA(s.longMA[:len(s.longMA)-1], s.longPeriod)

	signals := make([]strategy.Signal, 0)

	// 金叉:短期均线上穿长期均线,买入信号
	if prevShortAvg <= prevLongAvg && shortAvg > longAvg {
		signals = append(signals, strategy.Signal{
			StrategyID: s.name,
			Symbol:     s.extractSymbol(evt),
			Intent:     strategy.IntentLong,
			TargetQty:  100, // 买入100股
			Price:      price,
		})
	}

	// 死叉:短期均线下穿长期均线,卖出信号
	if prevShortAvg >= prevLongAvg && shortAvg < longAvg {
		signals = append(signals, strategy.Signal{
			StrategyID: s.name,
			Symbol:     s.extractSymbol(evt),
			Intent:     strategy.IntentFlat,
			TargetQty:  0, // 平仓
			Price:      price,
		})
	}

	return signals, nil
}

// OnStop 停止策略
func (s *MACrossStrategy) OnStop(ctx strategy.Context) error {
	// 清理资源
	s.shortMA = nil
	s.longMA = nil
	return nil
}

// calculateMA 计算移动平均线
func (s *MACrossStrategy) calculateMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	sum := 0.0
	start := len(prices) - period
	for i := start; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}

// extractPrice 从事件中提取价格(简化处理)
func (s *MACrossStrategy) extractPrice(evt market.Event) float64 {
	// TODO: 根据实际事件类型解析价格
	// 这里假设Data字段包含价格信息
	if bar, ok := evt.Data.(market.Bar); ok {
		return bar.Close
	}
	return 0
}

// extractSymbol 从事件中提取标的代码
func (s *MACrossStrategy) extractSymbol(evt market.Event) string {
	// TODO: 根据实际事件类型解析标的代码
	if bar, ok := evt.Data.(market.Bar); ok {
		return bar.Instrument.Symbol
	}
	return ""
}
