package execution

import (
	"context"
	"quant-trading/internal/domain/execution"
)

/*
Listener
========

Execution Engine 通过 Listener 向系统回传执行结果
*/
type Listener interface {
	OnExecutionEvent(ctx context.Context, event *execution.Event)
}
