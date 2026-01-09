package strategy

/*
最小定义
*/

import (
	"context"
	"quant-trading/internal/domain"
)

// Strategy 是所有策略的最小接口
// 核心原则：
// - 输入：行情 + Context
// - 输出：Signal（或 nil）
type Strategy interface {
	Name() string
	OnBar(ctx context.Context, sc Context, bar domain.MarketBar) *domain.Signal
}
