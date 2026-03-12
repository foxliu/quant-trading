package execution

import (
	"context"
	"quant-trading/internal/domain/execution"
)

/*
Engine
======

Execution Engine 的职责：
- 接收 Risk Engine 校验通过的 Order
- 将 Order 提交到真实 / 模拟 执行系统
- 不做任何业务判断
- 不阻塞上游业务流程

注意：
- Execution 不保证同步成交
- 成交结果通过 Execution Event 回传
*/
type Engine interface {
	// Submit 提交一个业务订单进行执行
	Submit(ctx context.Context, ord *execution.Order) error

	// RegisterListener 注册执行回报监听器
	RegisterListener(listener Listener)
}
