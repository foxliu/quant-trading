package strategyengine

import "quant-trading/internal/domain/strategy"

/*
Registry
========

策略注册表，用于在系统启动阶段完成策略装配。

注意：
- Registry 只在启动期使用
- 运行期不允许动态修改（避免状态混乱）
*/
type Registry struct {
	strategies []strategy.Strategy
}

func NewRegistry() *Registry {
	return &Registry{
		strategies: make([]strategy.Strategy, 0),
	}
}

func (r *Registry) Register(s strategy.Strategy) {
	r.strategies = append(r.strategies, s)
}

func (r *Registry) All() []strategy.Strategy {
	return r.strategies
}
