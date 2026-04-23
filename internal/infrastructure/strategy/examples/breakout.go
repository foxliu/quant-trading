package examples

import (
	"quant-trading/internal/domain/market"
	"quant-trading/internal/domain/strategy"
)

/*
BreakoutStrategy
================

突破策略示例。

策略逻辑:
- 价格突破N日最高价时买入
- 价格跌破N日最低价时卖出

参数:
- period: 回溯周期
- quantity: 交易数量
*/
type BreakoutStrategy struct {
	name string

	// 参数
	period   int
	quantity float64

	// 状态
	highs []float64 // 最高价序列
	lows  []float64 // 最低价序列
}

// NewBreakoutStrategy 创建突破策略
func NewBreakoutStrategy(name string, period int, quantity float64) *BreakoutStrategy {
	return &BreakoutStrategy{
		name:     name,
		period:   period,
		quantity: quantity,
		highs:    make([]float64, 0, period+1),
		lows:     make([]float64, 0, period+1),
	}
}

// Name 返回策略名称
func (s *BreakoutStrategy) Name() string {
	return s.name
}

// OnInit 初始化策略
func (s *BreakoutStrategy) OnInit(ctx strategy.Context) error {
	s.highs = make([]float64, 0, s.period+1)
	s.lows = make([]float64, 0, s.period+1)
	return nil
}

// OnMarketEvent 处理行情事件
func (s *BreakoutStrategy) OnMarketEvent(ctx strategy.Context, evt market.Event) ([]strategy.Signal, error) {
	// 提取价格
	high, low, close := s.extractPrices(evt)
	if high <= 0 || low <= 0 || close <= 0 {
		return nil, nil
	}

	// 更新价格序列
	s.highs = append(s.highs, high)
	s.lows = append(s.lows, low)

	// 保持序列长度
	if len(s.highs) > s.period+1 {
		s.highs = s.highs[1:]
		s.lows = s.lows[1:]
	}

	// 数据不足
	if len(s.highs) <= s.period {
		return nil, nil
	}

	// 计算前N日的最高价和最低价
	periodHigh := s.findMax(s.highs[:len(s.highs)-1])
	periodLow := s.findMin(s.lows[:len(s.lows)-1])

	signals := make([]strategy.Signal, 0)

	// 向上突破:当前价格突破前N日最高价
	if close > periodHigh {
		signals = append(signals, strategy.Signal{
			StrategyID: s.name,
			Symbol:     s.extractSymbol(evt),
			Intent:     strategy.IntentLong,
			TargetQty:  s.quantity,
			Price:      close,
		})
	}

	// 向下突破:当前价格跌破前N日最低价
	if close < periodLow {
		signals = append(signals, strategy.Signal{
			StrategyID: s.name,
			Symbol:     s.extractSymbol(evt),
			Intent:     strategy.IntentFlat,
			TargetQty:  0,
			Price:      close,
		})
	}

	return signals, nil
}

// OnStop 停止策略
func (s *BreakoutStrategy) OnStop(ctx strategy.Context) error {
	s.highs = nil
	s.lows = nil
	return nil
}

// findMax 找到序列中的最大值
func (s *BreakoutStrategy) findMax(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}

	max := prices[0]
	for _, p := range prices {
		if p > max {
			max = p
		}
	}
	return max
}

// findMin 找到序列中的最小值
func (s *BreakoutStrategy) findMin(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}

	min := prices[0]
	for _, p := range prices {
		if p < min {
			min = p
		}
	}
	return min
}

// extractPrices 从事件中提取价格
func (s *BreakoutStrategy) extractPrices(evt market.Event) (high, low, close float64) {
	if bar, ok := evt.Data.(market.Bar); ok {
		return bar.High, bar.Low, bar.Close
	}
	return 0, 0, 0
}

// extractSymbol 从事件中提取标的代码
func (s *BreakoutStrategy) extractSymbol(evt market.Event) string {
	if bar, ok := evt.Data.(market.Bar); ok {
		return bar.Instrument.Symbol
	}
	return ""
}
