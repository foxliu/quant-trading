package risk

import (
	"context"
)

/*
Engine
======

Risk Engine 的职责（V1）：

- 接收 Position Engine 输出的 Order
- 进行风控校验与裁剪
- 放行或拒绝
*/
type Engine interface {
	Start(ctx context.Context) error
	Stop() error

	Evaluate()
	Results() <-chan *Result
}
