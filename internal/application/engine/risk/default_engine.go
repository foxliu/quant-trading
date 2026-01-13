package risk

import (
	"context"
	"quant-trading/internal/domain/strategy"
)

/*
DefaultEngine
=============

Risk Engine 的默认实现：

处理链路：
Signal

	→ RiskCheck
	→ PositionManager
	→ ExecutionEngine
*/
type DefaultEngine struct {
	riskChecker RiskChecker
	positionMgr PositionManager
	executor    ExecutionEngine
}

func NewDefaultEngine(checker RiskChecker, pos PositionManager, exec ExecutionEngine) *DefaultEngine {
	return &DefaultEngine{
		riskChecker: checker,
		positionMgr: pos,
		executor:    exec,
	}
}

func (e *DefaultEngine) Start(ctx context.Context) error {
	return nil
}

func (e *DefaultEngine) Stop() error {
	return nil
}

func (e *DefaultEngine) Consume(signal strategy.Signal) {
	// 1. 风控检查
	if err := e.riskChecker.Check(signal); err != nil {
		return
	}

	// 2. 仓位处理
	if err := e.positionMgr.Apply(signal); err != nil {
		return
	}

	// 3. 下单执行
	_ = e.executor.Execute(signal)
}
